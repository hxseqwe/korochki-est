import React, { useState, useEffect } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import axios from 'axios';
import './Applications.css';

function Applications() {
  const [applications, setApplications] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const navigate = useNavigate();

  useEffect(() => {
    fetchApplications();
  }, []);

  const fetchApplications = async () => {
    try {
      const response = await axios.get('/api/applications');
      setApplications(response.data || []);
    } catch (err) {
      if (err.response?.status === 401) {
        localStorage.removeItem('user');
        navigate('/login');
      } else {
        setError('Ошибка загрузки заявок');
      }
    } finally {
      setLoading(false);
    }
  };

  const addReview = async (applicationId, review) => {
    try {
      await axios.post(`/api/applications/${applicationId}/review`, { review });
      fetchApplications();
    } catch (err) {
      setError('Ошибка при добавлении отзыва');
    }
  };

  const handleLogout = () => {
    localStorage.removeItem('user');
    navigate('/login');
  };

  let user = { fullName: '', login: '' };
  try {
    const userStr = localStorage.getItem('user');
    if (userStr) {
      user = JSON.parse(userStr);
    }
  } catch (e) {
    console.error('Error parsing user data');
  }

  if (loading) {
    return <div className="loading">Загрузка...</div>;
  }

  return (
    <div className="container">
      <header>
        <div>
          <h1>Корочки.есть</h1>
          <p className="welcome-message">Добро пожаловать, {user.fullName || user.login || 'пользователь'}!</p>
        </div>
        <div className="user-info">
          <button onClick={handleLogout} className="btn-link">Выйти</button>
        </div>
      </header>
      
      <nav className="nav-menu">
        <Link to="/" className="active">Мои заявки</Link>
        <Link to="/applications/new">Новая заявка</Link>
      </nav>
      
      <main>
        <h2>Мои заявки на обучение</h2>
        
        {error && <div className="error-message">{error}</div>}
        
        {applications && applications.length > 0 ? (
          <div className="applications-list">
            {applications.map(app => (
              <div key={app.id} className="application-card">
                <div className="application-header">
                  <h3>{app.course_name || 'Без названия'}</h3>
                  <span className={`status status-${app.status || 'new'}`}>
                    {app.status === 'new' && 'Новая'}
                    {app.status === 'in_progress' && 'Идет обучение'}
                    {app.status === 'completed' && 'Обучение завершено'}
                    {app.status === 'rejected' && 'Отклонена'}
                    {!app.status && 'Новая'}
                  </span>
                </div>
                
                <div className="application-details">
                  <p><strong>Дата начала:</strong> {app.start_date ? new Date(app.start_date).toLocaleDateString() : 'Не указана'}</p>
                  <p><strong>Способ оплаты:</strong> 
                    {app.payment_method === 'cash' ? ' Наличными' : app.payment_method === 'transfer' ? ' Перевод по номеру телефона' : ' Не указан'}
                  </p>
                  <p><strong>Дата подачи:</strong> {app.created_at ? new Date(app.created_at).toLocaleString() : 'Не указана'}</p>
                </div>
                
                {app.status === 'completed' && !app.review && (
                  <form onSubmit={(e) => {
                    e.preventDefault();
                    const review = e.target.review.value;
                    if (review && review.trim()) {
                      addReview(app.id, review);
                      e.target.reset();
                    }
                  }} className="review-form">
                    <textarea 
                      name="review" 
                      placeholder="Оставьте отзыв о качестве образовательных услуг" 
                      required 
                      rows="3"
                    />
                    <button type="submit" className="btn-secondary">Отправить отзыв</button>
                  </form>
                )}
                
                {app.review && (
                  <div className="review">
                    <p><strong>Ваш отзыв:</strong> {app.review}</p>
                  </div>
                )}
              </div>
            ))}
          </div>
        ) : (
          <p className="empty-state">
            У вас пока нет заявок. <Link to="/applications/new">Создать первую заявку</Link>
          </p>
        )}
      </main>
    </div>
  );
}

export default Applications;
