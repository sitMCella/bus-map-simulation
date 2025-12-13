import BusMap from './components/BusMap';
import { BrowserRouter, Route, Routes } from 'react-router-dom';
import './App.css';

function App() {
  return (
    <div>
      <BrowserRouter>
        <Routes>
          <Route path="/" element={<BusMap />} />
        </Routes>
      </BrowserRouter>
    </div>
  );
}

export default App;
