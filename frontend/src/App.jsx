import "./App.css";
import AppointmentPage from "./appointment";

function Navbar() {
  return (
    <nav className="navbar">
      <div className="brand">
        <div className="brand-icon">H</div>
        <span>H-Lab</span>
      </div>
    </nav>
  );
}

function App() {
  return (
    <div className="app">
      <Navbar />
      <AppointmentPage />
    </div>
  );
}

export default App;