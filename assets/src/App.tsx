import { useCallback, useEffect, useState } from 'react';
import './App.css';
import { MessageLog } from './MessageLog';
import { OnlineStatus } from './OnlineStatus';

function App() {

  const [messages, setMessages] = useState([])
  const [groups, setGroups] = useState([])

  const refresh = useCallback(async () => {
    fetch('/api/messages')
      .then(res => res.json())
      .then(({data}) => setMessages(data))
    fetch('/api/groups')
      .then(res => res.json())
      .then(({data}) => setGroups(data))
  }, [])

  useEffect(() => {
    const interval = setInterval(refresh, 2000)
    return () => clearInterval(interval)
  }, [refresh])

  return (
    <div className="App">
      <MessageLog messages={messages} />
      <OnlineStatus groups={groups} />
    </div>
  );
}

export default App;
