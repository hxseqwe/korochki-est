import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import axios from 'axios';
import './NewApplication.css';

function NewApplication() {
  const [form, setForm] = useState({
    course_name: '',
    start_date: '',
    payment_method: 'cash'
  });
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError('');

    try {
      await axios.post('/api/applications', form);
      navigate('/');
    } catch (err) {
      setError(err.response?.data || 'Ошибка при создании заявки');
    } finally {
      setLoading(false);
    }
  };

  const handleLogout = () => {
    localStorage.removeItem('user');
    navigate('/login');
  };

  const user = JSON.parse(localStorage.getItem('user') || '{}');
  const today = new Date().toISOString().split('T')[0];

  return (
    <div className="container">
      <header>
        <div>
          <h1>Корочки.есть</h1>
          <p className="welcome-message">Добро пожаловать, {user.fullName || user.login}!</p>
        </div>
        <div className="user-info">
          <button onClick={handleLogout} className="btn-link">Выйти</button>
        </div>
      </header>
      
      <nav className="nav-menu">
        <Link to="/">Мои заявки</Link>
        <Link to="/applications/new" className="active">Новая заявка</Link>
      </nav>
      
      <main>
        <h2>Новая заявка на обучение</h2>
        
        {error && <div className="error-message">{error}</div>}
        
        <form onSubmit={handleSubmit} className="application-form">
          <div className="form-group">
            <label>Наименование курса</label>
            <input
              type="text"
              value={form.course_name}
              onChange={(e) => setForm({...form, course_name: e.target.value})}
              required
            />
          </div>
          
          <div className="form-group">
            <label>Желаемая дата начала обучения</label>
            <input
              type="date"
              value={form.start_date}
              onChange={(e) => setForm({...form, start_date: e.target.value})}
              min={today}
              required
            />
          </div>
          
          <div className="form-group">
            <label>Способ оплаты</label>
            <div className="radio-group">
              <label>
                <input
                  type="radio"
                  value="cash"
                  checked={form.payment_method === 'cash'}
                  onChange={(e) => setForm({...form, payment_method: e.target.value})}
                />
                Наличными
              </label>
              <label>
                <input
                  type="radio"
                  value="transfer"
                  checked={form.payment_method === 'transfer'}
                  onChange={(e) => setForm({...form, payment_method: e.target.value})}
                />
                Перевод по номеру телефона
              </label>
            </div>
          </div>
          
          <button type="submit" className="btn-primary" disabled={loading}>
            {loading ? 'Отправка...' : 'Отправить заявку'}
          </button>
        </form>
      </main>
    </div>
  );
}

export default NewApplication;
