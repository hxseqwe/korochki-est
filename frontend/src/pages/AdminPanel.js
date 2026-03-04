import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import axios from 'axios';
import './AdminPanel.css';

function AdminPanel() {
  const [applications, setApplications] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [editingId, setEditingId] = useState(null);
  const [editForm, setEditForm] = useState({
    course_name: '',
    start_date: '',
    payment_method: ''
  });
  const navigate = useNavigate();

  useEffect(() => {
    fetchApplications();
  }, []);

  const fetchApplications = async () => {
    try {
      const response = await axios.get('/api/admin/applications');
      setApplications(response.data || []);
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

  const startEdit = (app) => {
    setEditingId(app.id);
    setEditForm({
      course_name: app.course_name,
      start_date: app.start_date.split('T')[0],
      payment_method: app.payment_method
    });
  };

  const cancelEdit = () => {
    setEditingId(null);
    setEditForm({ course_name: '', start_date: '', payment_method: '' });
  };

  const updateApplication = async (id) => {
    try {
      await axios.put(`/api/admin/applications/${id}`, editForm);
      setEditingId(null);
      fetchApplications();
    } catch (err) {
      setError('Ошибка при обновлении заявки');
    }
  };

  const deleteApplication = async (id) => {
    if (window.confirm('Вы уверены, что хотите удалить эту заявку?')) {
      try {
        await axios.delete(`/api/admin/applications/${id}`);
        fetchApplications();
      } catch (err) {
        setError('Ошибка при удалении заявки');
      }
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
  } catch (e) {}

  if (loading) {
    return <div className="loading">Загрузка...</div>;
  }

  const formatDate = (dateString) => {
    if (!dateString) return 'Не указана';
    return new Date(dateString).toLocaleDateString();
  };

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
        
        {applications && applications.length > 0 ? (
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
                    <td>
                      {editingId === app.id ? (
                        <input
                          type="text"
                          value={editForm.course_name}
                          onChange={(e) => setEditForm({...editForm, course_name: e.target.value})}
                          className="edit-input"
                        />
                      ) : (
                        app.course_name
                      )}
                    </td>
                    <td>
                      {editingId === app.id ? (
                        <input
                          type="date"
                          value={editForm.start_date}
                          onChange={(e) => setEditForm({...editForm, start_date: e.target.value})}
                          className="edit-input"
                        />
                      ) : (
                        formatDate(app.start_date)
                      )}
                    </td>
                    <td>
                      {editingId === app.id ? (
                        <select
                          value={editForm.payment_method}
                          onChange={(e) => setEditForm({...editForm, payment_method: e.target.value})}
                          className="edit-select"
                        >
                          <option value="cash">Наличными</option>
                          <option value="transfer">Перевод</option>
                        </select>
                      ) : (
                        app.payment_method === 'cash' ? 'Наличными' : 'Перевод'
                      )}
                    </td>
                    <td>
                      <select 
                        onChange={(e) => updateStatus(app.id, e.target.value)}
                        value={app.status || 'new'}
                        className="status-select"
                      >
                        <option value="new">Новая</option>
                        <option value="in_progress">Идет обучение</option>
                        <option value="completed">Обучение завершено</option>
                        <option value="rejected">Отклонена</option>
                      </select>
                    </td>
                    <td>
                      <div className="action-buttons">
                        {editingId === app.id ? (
                          <>
                            <button onClick={() => updateApplication(app.id)} className="btn-save">Сохранить</button>
                            <button onClick={cancelEdit} className="btn-cancel">Отмена</button>
                          </>
                        ) : (
                          <>
                            <button onClick={() => startEdit(app)} className="btn-edit">Редактировать</button>
                            <button onClick={() => deleteApplication(app.id)} className="btn-delete">Удалить</button>
                          </>
                        )}
                      </div>
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