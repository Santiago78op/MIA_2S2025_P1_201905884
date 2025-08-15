import React from 'react';
import ErrorNotification, { NotificationProps } from './ErrorNotification';

interface NotificationContainerProps {
  notifications: NotificationProps[];
  onRemove: (id: string) => void;
}

const NotificationContainer: React.FC<NotificationContainerProps> = ({
  notifications,
  onRemove
}) => {
  return (
    <div className="notification-container">
      {notifications.map((notification) => (
        <ErrorNotification
          key={notification.id}
          {...notification}
          onClose={onRemove}
        />
      ))}
    </div>
  );
};

export default NotificationContainer;