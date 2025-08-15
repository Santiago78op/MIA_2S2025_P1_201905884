import React, { useEffect, useState } from 'react';
import './ErrorNotification.css';

export interface NotificationProps {
  id: string;
  type: 'success' | 'warning' | 'error' | 'info';
  title: string;
  message: string;
  duration?: number;
  persistent?: boolean;
  actions?: Array<{
    label: string;
    action: () => void;
    style?: 'primary' | 'secondary';
  }>;
}

interface ErrorNotificationProps extends NotificationProps {
  onClose: (id: string) => void;
}

const ErrorNotification: React.FC<ErrorNotificationProps> = ({
  id,
  type,
  title,
  message,
  duration = 5000,
  persistent = false,
  actions = [],
  onClose
}) => {
  const [isVisible, setIsVisible] = useState(true);
  const [isAnimating, setIsAnimating] = useState(false);

  useEffect(() => {
    setIsAnimating(true);
    
    if (!persistent && duration > 0) {
      const timer = setTimeout(() => {
        handleClose();
      }, duration);
      
      return () => clearTimeout(timer);
    }
  }, [duration, persistent]);

  const handleClose = () => {
    setIsAnimating(false);
    setTimeout(() => {
      setIsVisible(false);
      onClose(id);
    }, 300);
  };

  const getIcon = () => {
    switch (type) {
      case 'success': return '✅';
      case 'warning': return '⚠️';
      case 'error': return '❌';
      case 'info': return 'ℹ️';
      default: return '📝';
    }
  };

  const getTypeClass = () => {
    switch (type) {
      case 'success': return 'notification-success';
      case 'warning': return 'notification-warning';
      case 'error': return 'notification-error';
      case 'info': return 'notification-info';
      default: return 'notification-info';
    }
  };

  if (!isVisible) return null;

  return (
    <div 
      className={`error-notification ${getTypeClass()} ${isAnimating ? 'animate-in' : 'animate-out'}`}
      role="alert"
    >
      <div className="notification-header">
        <span className="notification-icon">{getIcon()}</span>
        <div className="notification-content">
          <h4 className="notification-title">{title}</h4>
          <p className="notification-message">{message}</p>
        </div>
        <button 
          className="notification-close"
          onClick={handleClose}
          aria-label="Cerrar notificación"
        >
          ✕
        </button>
      </div>
      
      {actions.length > 0 && (
        <div className="notification-actions">
          {actions.map((action, index) => (
            <button
              key={index}
              className={`notification-action ${action.style === 'primary' ? 'primary' : 'secondary'}`}
              onClick={() => {
                action.action();
                if (!persistent) handleClose();
              }}
            >
              {action.label}
            </button>
          ))}
        </div>
      )}
      
      {!persistent && duration > 0 && (
        <div className="notification-progress">
          <div 
            className="progress-bar" 
            style={{ animationDuration: `${duration}ms` }}
          />
        </div>
      )}
    </div>
  );
};

export default ErrorNotification;