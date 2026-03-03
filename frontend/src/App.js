import React from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import Login from './pages/Login';
import Register from './pages/Register';
import Applications from './pages/Applications';
import NewApplication from './pages/NewApplication';
import AdminPanel from './pages/AdminPanel';
import './App.css';

function App() {
  const isAuthenticated = () => {
    return localStorage.getItem('user') !== null;
  };

  const isAdmin = () => {
    const user = JSON.parse(localStorage.getItem('user') || '{}');
    return user.isAdmin === true;
  };

  const PrivateRoute = ({ children }) => {
    return isAuthenticated() ? children : <Navigate to="/login" />;
  };

  const AdminRoute = ({ children }) => {
    return isAuthenticated() && isAdmin() ? children : <Navigate to="/" />;
  };

  return (
    <Router>
      <div className="App">
        <Routes>
          <Route path="/login" element={<Login />} />
          <Route path="/register" element={<Register />} />
          <Route path="/" element={
            <PrivateRoute>
              <Applications />
            </PrivateRoute>
          } />
          <Route path="/applications/new" element={
            <PrivateRoute>
              <NewApplication />
            </PrivateRoute>
          } />
          <Route path="/admin" element={
            <AdminRoute>
              <AdminPanel />
            </AdminRoute>
          } />
        </Routes>
      </div>
    </Router>
  );
}

export default App;
