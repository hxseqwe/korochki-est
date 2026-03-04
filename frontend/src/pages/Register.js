import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import axios from 'axios';
import './Auth.css';

function Register() {
  const [form, setForm] = useState({
    login: '',
    password: '',
    full_name: '',
    phone: '',
    email: ''
  });
  const [errors, setErrors] = useState({});
  const [serverError, setServerError] = useState('');
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

  const validateForm = () => {
    const newErrors = {};

    if (!form.login) {
      newErrors.login = 'Логин обязателен';
    } else if (form.login.length < 6) {
      newErrors.login = 'Логин должен содержать минимум 6 символов';
    } else if (!/^[A-Za-z0-9]+$/.test(form.login)) {
      newErrors.login = 'Логин может содержать только латиницу и цифры';
    }

    if (!form.password) {
      newErrors.password = 'Пароль обязателен';
    } else if (form.password.length < 8) {
      newErrors.password = 'Пароль должен содержать минимум 8 символов';
    } else if (!/[A-Z]/.test(form.password)) {
      newErrors.password = 'Пароль должен содержать хотя бы одну заглавную букву';
    } else if (!/[a-z]/.test(form.password)) {
      newErrors.password = 'Пароль должен содержать хотя бы одну строчную букву';
    } else if (!/[0-9]/.test(form.password)) {
      newErrors.password = 'Пароль должен содержать хотя бы одну цифру';
    } else if (!/[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]/.test(form.password)) {
      newErrors.password = 'Пароль должен содержать хотя бы один специальный символ';
    }

    if (!form.full_name) {
      newErrors.full_name = 'ФИО обязательно';
    } else if (!/^[А-Яа-яЁё\s]+$/.test(form.full_name)) {
      newErrors.full_name = 'ФИО может содержать только кириллицу и пробелы';
    }

    if (!form.phone) {
      newErrors.phone = 'Телефон обязателен';
    } else {
      const cleanedPhone = form.phone.replace(/\D/g, '');
      if (cleanedPhone.length !== 11) {
        newErrors.phone = 'Телефон должен содержать 11 цифр';
      } else if (!cleanedPhone.startsWith('8')) {
        newErrors.phone = 'Телефон должен начинаться с 8';
      } else {
        const formattedPhone = `8(${cleanedPhone.slice(1,4)})${cleanedPhone.slice(4,7)}-${cleanedPhone.slice(7,9)}-${cleanedPhone.slice(9,11)}`;
        setForm(prev => ({ ...prev, phone: formattedPhone }));
      }
    }

    if (!form.email) {
      newErrors.email = 'Email обязателен';
    } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(form.email)) {
      newErrors.email = 'Введите корректный email';
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handlePhoneChange = (e) => {
    let value = e.target.value.replace(/\D/g, '');
    
    if (value.length > 11) {
      value = value.slice(0, 11);
    }
    
    if (value.length > 0) {
      if (value.length <= 1) {
        value = `8(${value.slice(1)}`;
      } else if (value.length <= 4) {
        value = `8(${value.slice(1,4)}`;
      } else if (value.length <= 7) {
        value = `8(${value.slice(1,4)})${value.slice(4,7)}`;
      } else if (value.length <= 9) {
        value = `8(${value.slice(1,4)})${value.slice(4,7)}-${value.slice(7,9)}`;
      } else {
        value = `8(${value.slice(1,4)})${value.slice(4,7)}-${value.slice(7,9)}-${value.slice(9,11)}`;
      }
    }
    
    setForm({ ...form, phone: value });
    
    if (errors.phone) {
      setErrors({ ...errors, phone: '' });
    }
  };

  const handleChange = (e) => {
    const { name, value } = e.target;
    setForm({ ...form, [name]: value });
    
    if (errors[name]) {
      setErrors({ ...errors, [name]: '' });
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    
    if (!validateForm()) {
      return;
    }

    setLoading(true);
    setServerError('');

    try {
      const response = await axios.post('/api/register', form);
      localStorage.setItem('user', JSON.stringify(response.data));
      navigate('/');
    } catch (err) {
      setServerError(err.response?.data || 'Ошибка регистрации');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="auth-container">
      <div className="auth-card">
        <h1>Регистрация</h1>
        <h2>Корочки.есть</h2>
        
        {serverError && <div className="error-message">{serverError}</div>}
        
        <form onSubmit={handleSubmit} noValidate>
          <div className="form-group">
            <label>Логин (латиница и цифры, мин. 6 символов)</label>
            <input
              type="text"
              name="login"
              value={form.login}
              onChange={handleChange}
              className={errors.login ? 'error' : ''}
              placeholder="qwe123"
              required
            />
            {errors.login && <div className="field-error">{errors.login}</div>}
          </div>
          
          <div className="form-group">
            <label>Пароль (мин. 8 символов, заглавная, строчная, цифра, спецсимвол)</label>
            <input
              type="password"
              name="password"
              value={form.password}
              onChange={handleChange}
              className={errors.password ? 'error' : ''}
              required
            />
            {errors.password && <div className="field-error">{errors.password}</div>}
          </div>
          
          <div className="form-group">
            <label>ФИО (только кириллица и пробелы)</label>
            <input
              type="text"
              name="full_name"
              value={form.full_name}
              onChange={handleChange}
              className={errors.full_name ? 'error' : ''}
              placeholder="Иванов Иван Иванович"
              required
            />
            {errors.full_name && <div className="field-error">{errors.full_name}</div>}
          </div>
          
          <div className="form-group">
            <label>Телефон (формат: 8(XXX)XXX-XX-XX)</label>
            <input
              type="tel"
              name="phone"
              value={form.phone}
              onChange={handlePhoneChange}
              className={errors.phone ? 'error' : ''}
              placeholder="8(999)999-99-99"
              required
            />
            {errors.phone && <div className="field-error">{errors.phone}</div>}
          </div>
          
          <div className="form-group">
            <label>Email</label>
            <input
              type="email"
              name="email"
              value={form.email}
              onChange={handleChange}
              className={errors.email ? 'error' : ''}
              placeholder="example@mail.ru"
              required
            />
            {errors.email && <div className="field-error">{errors.email}</div>}
          </div>
          
          <button type="submit" className="btn-primary" disabled={loading}>
            {loading ? 'Регистрация...' : 'Создать пользователя'}
          </button>
        </form>
        
        <p className="auth-link">
          Уже зарегистрированы? <Link to="/login">Войти</Link>
        </p>
      </div>
    </div>
  );
}

export default Register;