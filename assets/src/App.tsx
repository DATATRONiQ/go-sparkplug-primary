import { useEffect, useState } from 'react';
import './App.css';
import { MessageLog } from './MessageLog';

function App() {

  const [messages, setMessages] = useState([])

  useEffect(() => {
    fetch('/api/messages')
      .then(res => res.json())
      .then(({data}) => setMessages(data))
  }, [])

  return (
    <div className="App">
      <MessageLog messages={messages} />
    </div>
  );
}

export default App;
