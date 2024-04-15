import React, { useState, useEffect } from 'react';
import './App.css'; // Stiller için ekstra CSS dosyası

function App() {
  const [ws, setWs] = useState(null);
  const [message, setMessage] = useState('');
  const [receivedMessages, setReceivedMessages] = useState([]);

  useEffect(() => {
    const newWs = new WebSocket('ws://localhost:8080/ws');
    newWs.onmessage = (event) => {
      setReceivedMessages(prev => [...prev, event.data]);
    };
    newWs.onopen = () => {
      console.log('WebSocket Connected');
    };
    newWs.onerror = (error) => {
      console.error('WebSocket Error:', error);
    };
    setWs(newWs);
    return () => {
      newWs.close();
    };
  }, []);

  const sendMessage = () => {
    if (ws && message !== '') {
      ws.send(message);
      setMessage('');
    }
  };

  return (
    <div className="app">
      <header className="app-header">
        Chat App
      </header>
      <div className="message-container">
        {receivedMessages.map((msg, index) => (
          <div key={index} className="message-bubble">
            {msg}
          </div>
        ))}
      </div>
      <div className="input-container">
        <input
          type="text"
          className="message-input"
          placeholder="Type a message..."
          value={message}
          onChange={(e) => setMessage(e.target.value)}
          onKeyPress={(e) => e.key === 'Enter' && sendMessage()}
        />
        <button className="send-button" onClick={sendMessage}>
          Send
        </button>
      </div>
    </div>
  );
}

export default App;
