// Package database provides the services for storing and retrieving data from the database.
package database

import (
   "database/sql"
   "fmt"
   "github.com/joho/godotenv"
   _ "github.com/lib/pq"
   "os"
   "strconv"
   "time"
)

// DatabaseConnection implements the PostgreSQL client.
type DatabaseConnection struct {
	Db *sql.DB
}

type BusStop struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Latitude string `json:"latitude"`
    Longitude string `json:"longitude"`
}

type Bus struct {
	Id string `json:"id"`
	Latitude string `json:"latitude"`
    Longitude string `json:"longitude"`
}

type BusTimeTable struct {
	BusId string `json:"bus_id"`
	BusStopId string `json:"bus_stop_id"`
	TimeSeconds time.Duration `json:"time_seconds"`
	Timestamp time.Time `json:"timestamp"`
}

type BusPosition struct {
	Id string `json:"id"`
	CreationTime time.Time `json:"creationtime"`
    BusId     string `json:"bus_id"`
    Latitude string `json:"latitude"`
    Longitude string `json:"longitude"`
	NextBusStopId string `json:"next_bus_stop_id"`
	IsBusStop bool `json:"is_bus_stop"`
}

// NewDatabaseConnection creates a new connection to PostgreSQL.
func NewDatabaseConnection() (databaseConnection DatabaseConnection, err error) {
	_ = godotenv.Load()
	host := os.Getenv("DB_HOST")
	port, _ := strconv.Atoi(os.Getenv("DB_PORT"))
	user := os.Getenv("DB_USER")
	dbname := os.Getenv("DB_NAME")
	pass := os.Getenv("DB_PASSWORD")
	
	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", user, pass, host, port, dbname)

	var db *sql.DB
	for i := 0; i < 10; i++ {
		db, err = sql.Open("postgres", dbUrl)
		if err == nil {
			err = db.Ping()
		}
		if err == nil {
			break
		}
		fmt.Println("Database not ready, retrying in 2s... (%d/10)", i+1)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		fmt.Println("Cannot connect to database:", err)
	}
	db.SetMaxOpenConns(80)
	db.SetMaxIdleConns(40)
	db.SetConnMaxLifetime(time.Duration(60) * time.Minute)
	databaseConnection = DatabaseConnection{
		Db: db,
	}
	return
}

// InitDatabase creates the tables in the Database if these don't exist.
func (dc DatabaseConnection) InitDatabase() error {
	err := dc.createBusStopTable()
	if err != nil {
		return err
	}
	err = dc.createBusStopEntries()
	if err != nil {
		return err
	}
	err = dc.createBusTable()
	if err != nil {
		return err
	}
	err = dc.createBusEntries()
	if err != nil {
		return err
	}
	err = dc.createBusTimeTable()
	if err != nil {
		return err
	}
	err = dc.createBusTimeEntries()
	if err != nil {
		return err
	}
	err = dc.createBusPositionTable()
	if err != nil {
		return err
	}
	err = dc.createFunction()
	if err != nil {
		return err
	}
	err = dc.dropTrigger()
	if err != nil {
		return err
	}
	return dc.createTrigger()
}

// Close closes the Database connection.
func (dc DatabaseConnection) Close() error {
	return dc.Db.Close()
}

func (dc DatabaseConnection) executeTransaction(sqlStmt string) (err error) {
	tx, err := dc.Db.Begin()
	if err != nil {
		return
	}
	defer func() {
		switch err {
		case nil:
			sqlerr := tx.Commit()
			if err == nil {
				err = sqlerr
			}
		default:
			sqlerr := tx.Rollback()
			if err == nil {
				err = sqlerr
			}
		}
	}()
	_, err = dc.Db.Exec(sqlStmt)
	return
}

func (dc DatabaseConnection) createBusStopTable() (err error) {
	sqlStmt := `CREATE TABLE IF NOT EXISTS bus_stop
				(
					id varchar (36) NOT NULL,
					name varchar (36) NOT NULL,
					latitude DOUBLE PRECISION NOT NULL,
					longitude DOUBLE PRECISION NOT NULL,
					PRIMARY KEY(id)
				);`
	err = dc.executeTransaction(sqlStmt)
	return
}

func (dc DatabaseConnection) createBusStopEntries() (err error) {
	sqlStmt := `INSERT INTO bus_stop (id, name, latitude, longitude) VALUES 
					('1', 'Stazione Tiburtina', 41.9096, 12.52975),
					('2', 'Tiburtina / Crociate', 41.90815, 12.52589),
					('3', 'Tiburtina / Valerio Massimo', 41.90594, 12.52228),
					('4', 'Tiburtina / Castro Laurenziano', 41.90391, 12.52032),
					('5', 'De Lollis / Verano', 41.90171, 12.5178),
					('6', 'De Lollis / Irpini', 41.90101, 12.51577),
					('7', 'Ramni / Marrucini', 41.9001, 12.51309),
					('8', 'Ramni/ Porta Tiburtina', 41.89906, 12.51032),
					('9', 'Pretoriano', 41.90064, 12.50818),
					('10', 'Catro Pretorio / Monzambano', 41.90347, 12.50687),
					('11', 'S. M. Battaglia', 41.90604, 12.50557),
					('12', 'Indipendenza', 41.90475, 12.50249),
					('13', 'Volturno / Gaeta', 41.90404, 12.50032),
					('14', 'Volturno / Cernaia', 41.90521, 12.49892),
					('15', 'Palestro', 41.90797, 12.50053),
					('16', 'XX Settembre / Piave', 41.90709, 12.49814),
					('17', 'XX Settembre / Min. Finanze', 41.90592, 12.49644),
					('18', 'Bissolati', 41.90525, 12.49266),
					('19', 'Barberini', 41.9043, 12.48889),
					('20', 'Tritone / Berberini', 41.90343, 12.48755),
					('21', 'Tritone / Fontana Trevi', 41.90262, 12.48446),
					('22', 'S. Claudio', 41.90195, 12.48037),
					('23', 'Corso / Minghetti', 41.89945, 12.48107),
					('24', 'Plebiscito', 41.89634, 12.48062),
					('25', 'Argentina', 41.89608, 12.47684),
					('26', 'Rinascimento', 41.89814, 12.47398),
					('27', 'Senato', 41.90028, 12.47382),
					('28', 'Zanardelli', 41.90139, 12.47211),
					('29', 'Lungotevere Marzio', 41.90316, 12.47378),
					('30', 'Vittoria Colonna', 41.90517, 12.47168),
					('31', 'Piazza Cavour', 41.90588, 12.46994),
					('32', 'Crescenzo / Orazio', 41.90547, 12.46739),
					('33', 'Crescenzo / Terenzio', 41.90572, 12.46387),
					('34', 'Crescenzo / Rinascimento', 41.90605, 12.45911),
					('35', 'Bastioni di Michelangelo', 41.90694, 12.45573),
					('36', 'Leone IV', 41.90903, 12.45524),
					('37', 'Doria A. / Largo Trionfale', 41.91007, 12.45347),
					('38', 'Di Lauria', 41.90875, 12.4503),
					('39', 'Emo', 41.9069, 12.44926),
					('40', 'Stazione Metro Cipro', 41.90722, 12.44789)
					ON CONFLICT (id) DO NOTHING;`
	err = dc.executeTransaction(sqlStmt)
	return
}

func (dc DatabaseConnection) createBusTable() (err error) {
	sqlStmt := `CREATE TABLE IF NOT EXISTS bus
				(
					id varchar (36) NOT NULL,
					latitude DOUBLE PRECISION NOT NULL,
					longitude DOUBLE PRECISION NOT NULL,
					PRIMARY KEY(id)
				);`
	err = dc.executeTransaction(sqlStmt)
	return
}

func (dc DatabaseConnection) createBusEntries() (err error) {
	sqlStmt := `INSERT INTO bus(id, latitude, longitude) VALUES
					('492', 41.9096, 12.52975)
					ON CONFLICT (id) DO NOTHING;`
	err = dc.executeTransaction(sqlStmt)
	return
}

func (dc DatabaseConnection) createBusTimeTable() (err error) {
	sqlStmt := `CREATE TABLE IF NOT EXISTS bus_time_table
				(
					bus_id varchar (36) NOT NULL REFERENCES bus(id),
					bus_stop_id varchar (36) NOT NULL REFERENCES bus_stop(id),
					time_seconds INTEGER NOT NULL,
					PRIMARY KEY(bus_id, bus_stop_id)
				);`
	err = dc.executeTransaction(sqlStmt)
	return
}

func (dc DatabaseConnection) createBusTimeEntries() (err error) {
	sqlStmt := `INSERT INTO bus_time_table(bus_id, bus_stop_id, time_seconds) VALUES
					('492', '1', 0),
					('492', '2', 51),
					('492', '3', 70),
					('492', '4', 78),
					('492', '5', 110),
					('492', '6', 117),
					('492', '7', 127),
					('492', '8', 147),
					('492', '9', 170),
					('492', '10', 196),
					('492', '11', 228),
					('492', '12', 262),
					('492', '13', 288),
					('492', '14', 303),
					('492', '15', 330),
					('492', '16', 339),
					('492', '17', 346),
					('492', '18', 380),
					('492', '19', 406),
					('492', '20', 427),
					('492', '21', 464),
					('492', '22', 486),
					('492', '23', 507),
					('492', '24', 529),
					('492', '25', 551),
					('492', '26', 591),
					('492', '27', 620),
					('492', '28', 684),
					('492', '29', 715),
					('492', '30', 750),
					('492', '31', 783),
					('492', '32', 803),
					('492', '33', 814),
					('492', '34', 829),
					('492', '35', 879),
					('492', '36', 891),
					('492', '37', 906),
					('492', '38', 923),
					('492', '39', 952),
					('492', '40', 969)
					ON CONFLICT (bus_id, bus_stop_id) DO NOTHING;`
	err = dc.executeTransaction(sqlStmt)
	return
}

func (dc DatabaseConnection) createBusPositionTable() (err error) {
	sqlStmt := `CREATE TABLE IF NOT EXISTS bus_position
				(
					id bigserial NOT NULL,
					creationtime timestamp NOT NULL DEFAULT NOW(),
					bus_id varchar (36) NOT NULL REFERENCES bus(id),
					latitude DOUBLE PRECISION NOT NULL,
					longitude DOUBLE PRECISION NOT NULL,
					next_bus_stop_id varchar (36) NOT NULL REFERENCES bus_stop(id),
					is_bus_stop bool NOT NULL,
					PRIMARY KEY(id)
				);`
	err = dc.executeTransaction(sqlStmt)
	return
}

func (dc DatabaseConnection) createFunction() (err error) {
	sqlStmt := `CREATE OR REPLACE FUNCTION notify_bus_position_event() RETURNS TRIGGER AS
				$$
				BEGIN
					PERFORM pg_notify('bus_position_notification', json_build_object(
					'id', NEW.id,
					'creationtime', to_char(NEW.creationtime, 'YYYY-MM-DD"T"HH24:MI:SS.US"Z"'),
					'busId', NEW.bus_id,
					'latitude', NEW.latitude,
					'longitude', NEW.longitude,
					'nextBusStopId', NEW.next_bus_stop_id,
					'isBusStop', NEW.is_bus_stop
					)::text);
					RETURN NULL;
				END;
				$$
				LANGUAGE plpgsql;;`
	err = dc.executeTransaction(sqlStmt)
	return
}

func (dc DatabaseConnection) dropTrigger() (err error) {
	sqlStmt := `DROP TRIGGER IF EXISTS notify_bus_position ON bus_position;;`
	err = dc.executeTransaction(sqlStmt)
	return
}

func (dc DatabaseConnection) createTrigger() (err error) {
	sqlStmt := `CREATE TRIGGER notify_bus_position
					AFTER INSERT
					ON bus_position
					FOR EACH ROW
				EXECUTE PROCEDURE notify_bus_position_event();;`
	err = dc.executeTransaction(sqlStmt)
	return
}

func (dc DatabaseConnection) GetBusStopEntries() (error, []BusStop) {
	sqlStmt, err := dc.Db.Prepare("SELECT id, name, latitude, longitude FROM bus_stop")
	if err != nil {
		return err, nil
    }
    defer sqlStmt.Close()

	var busStopEntries []BusStop

	rows, err := sqlStmt.Query()
	if err != nil {
		return err, nil
    }

	for rows.Next() {
		var bs BusStop
		err := rows.Scan(
			&bs.Id,
			&bs.Name,
			&bs.Latitude,
			&bs.Longitude,
		)
		if err != nil {
			return err, nil
		}
		busStopEntries = append(busStopEntries, bs)
	}
	if rows.Err() != nil {
		return rows.Err(), nil
	}
	return nil, busStopEntries
}

func (dc DatabaseConnection) GetBusEntries() (error, []Bus) {
	sqlStmt, err := dc.Db.Prepare("SELECT id, latitude, longitude FROM bus")
	if err != nil {
		return err, nil
    }
    defer sqlStmt.Close()

	var busEntries []Bus

	rows, err := sqlStmt.Query()
	if err != nil {
		return err, nil
    }

	for rows.Next() {
		var b Bus
		err := rows.Scan(
			&b.Id,
			&b.Latitude,
			&b.Longitude,
		)
		if err != nil {
			return err, nil
		}
		busEntries = append(busEntries, b)
	}
	if rows.Err() != nil {
		return rows.Err(), nil
	}
	return nil, busEntries
}

func (dc DatabaseConnection) GetBusTimeTableEntries(busId string) (error, []BusTimeTable) {
	sqlStmt, err := dc.Db.Prepare("SELECT bus_id, bus_stop_id, time_seconds FROM bus_time_table WHERE bus_id LIKE $1")
	if err != nil {
		return err, nil
    }
    defer sqlStmt.Close()
	

	var busTimeTableEntries []BusTimeTable

	rows, err := sqlStmt.Query(busId)
	if err != nil {
		return err, nil
    }

	var currentTime = time.Now().Local()

	for rows.Next() {
		fmt.Println("scam row")
		var btt BusTimeTable
		err := rows.Scan(
			&btt.BusId,
			&btt.BusStopId,
			&btt.TimeSeconds,
		)
		if err != nil {
			return err, nil
		}
		btt.Timestamp = currentTime.Add(time.Second * btt.TimeSeconds)
		busTimeTableEntries = append(busTimeTableEntries, btt)
	}
	if rows.Err() != nil {
		return rows.Err(), nil
	}
	return nil, busTimeTableEntries
}

func (dc DatabaseConnection) BusExists(busId string) (err error, exists bool) {
	sqlStmt, err := dc.Db.Prepare("SELECT COUNT(*) FROM bus WHERE id LIKE $1")
	if err != nil {
		return
    }
    defer sqlStmt.Close()
	var count int
    err = sqlStmt.QueryRow(busId).Scan(&count)
	if err != nil {
		if err != sql.ErrNoRows {
			return
		}
		return nil, false
	}
	exists = false
	if count == 1 {
		exists = true
	}
	return
}

func (dc DatabaseConnection) CreateBus(busId string, latitude string, longitude string) (err error) {
	sqlStmt, err := dc.Db.Prepare("INSERT INTO bus (id, latitude, longitude) VALUES ($1, $2, $3)")
	if err != nil {
		return
    }
    defer sqlStmt.Close()
	_, err = sqlStmt.Exec(busId, latitude, longitude)
	return
}

func (dc DatabaseConnection) CreateBusPosition(busId string, latitude string, longitude string, nextBusStopId string, isBusStop bool) (err error, bp BusPosition) {
	sqlStmt, err := dc.Db.Prepare("INSERT INTO bus_position (bus_id, latitude, longitude, next_bus_stop_id, is_bus_stop) VALUES ($1, $2, $3, $4, $5) RETURNING id, creationtime, bus_id, latitude, longitude, next_bus_stop_id, is_bus_stop")
	if err != nil {
		return
    }
    defer sqlStmt.Close()
	err = sqlStmt.QueryRow(busId, latitude, longitude, nextBusStopId, isBusStop).Scan(
		&bp.Id,
		&bp.CreationTime,
		&bp.BusId,
		&bp.Latitude,
		&bp.Longitude,
		&bp.NextBusStopId,
		&bp.IsBusStop,
	)
	return
}
