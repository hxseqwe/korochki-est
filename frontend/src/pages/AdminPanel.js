import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import axios from 'axios';
import './AdminPanel.css';

function AdminPanel() {
  const [applications, setApplications] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const navigate = useNavigate();

  useEffect(() => {
    fetchApplications();
  }, []);

  const fetchApplications = async () => {
    try {
      const response = await axios.get('/api/admin/applications');
      setApplications(response.data);
    } catch (err) {
      if (err.response?.status === 401 || err.response?.status === 403) {
        localStorage.removeItem('user');
        navigate('/login');
      } else {
        setError('Ошибка загрузки заявок');
      }
    } finally {
      setLoading(false);
    }
  };

  const updateStatus = async (applicationId, status) => {
    try {
      await axios.post(`/api/admin/applications/${applicationId}/status`, { status });
      fetchApplications();
    } catch (err) {
      setError('Ошибка при обновлении статуса');
    }
  };

  const handleLogout = () => {
    localStorage.removeItem('user');
    navigate('/login');
  };

  const user = JSON.parse(localStorage.getItem('user') || '{}');

  if (loading) return <div className="loading">Загрузка...</div>;

  return (
    <div className="container">
      <header>
        <div>
          <h1>Корочки.есть - Админ-панель</h1>
          <p className="welcome-message">Администратор: {user.fullName || user.login}</p>
        </div>
        <div className="user-info">
          <button onClick={handleLogout} className="btn-link">Выйти</button>
        </div>
      </header>
      
      <main>
        <h2>Управление заявками</h2>
        
        {error && <div className="error-message">{error}</div>}
        
        {applications.length > 0 ? (
          <table className="applications-table">
            <thead>
              <tr>
                <th>ID</th>
                <th>Пользователь</th>
                <th>Курс</th>
                <th>Дата начала</th>
                <th>Способ оплаты</th>
                <th>Статус</th>
                <th>Действия</th>
              </tr>
            </thead>
            <tbody>
              {applications.map(app => (
                <React.Fragment key={app.id}>
                  <tr>
                    <td>{app.id}</td>
                    <td>
                      {app.user?.full_name}<br/>
                      <small>{app.user?.login}</small>
                    </td>
                    <td>{app.course_name}</td>
                    <td>{new Date(app.start_date).toLocaleDateString()}</td>
                    <td>{app.payment_method === 'cash' ? 'Наличными' : 'Перевод'}</td>
                    <td>
                      <span className={`status status-${app.status}`}>
                        {app.status === 'new' && 'Новая'}
                        {app.status === 'in_progress' && 'Идет обучение'}
                        {app.status === 'completed' && 'Завершено'}
                        {app.status === 'rejected' && 'Отклонена'}
                      </span>
                    </td>
                    <td>
                      <select 
                        onChange={(e) => updateStatus(app.id, e.target.value)}
                        value={app.status}
                        className="status-select"
                      >
                        <option value="new">Новая</option>
                        <option value="in_progress">Идет обучение</option>
                        <option value="completed">Обучение завершено</option>
                        <option value="rejected">Отклонена</option>
                      </select>
                    </td>
                  </tr>
                  {app.review && (
                    <tr className="review-row">
                      <td colSpan="7">
                        <strong>Отзыв:</strong> {app.review}
                      </td>
                    </tr>
                  )}
                </React.Fragment>
              ))}
            </tbody>
          </table>
        ) : (
          <p className="empty-state">Нет заявок для отображения</p>
        )}
      </main>
    </div>
  );
}

export default AdminPanel;
