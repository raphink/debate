import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import ErrorBoundary from './components/common/ErrorBoundary/ErrorBoundary';
import Home from './pages/Home';
import PanelistSelection from './pages/PanelistSelection';
import DebateGeneration from './pages/DebateGeneration';
import DebateViewer from './pages/DebateViewer';
import DebateHistory from './pages/DebateHistory';
import NotFound from './pages/NotFound';
import './App.css';

function App() {
  return (
    <ErrorBoundary>
      <Router>
        <div className="App">
          <header className="App-header">
            <h1>Theology & Philosophy Debate Generator</h1>
          </header>
          
          <main className="App-main">
            <Routes>
              <Route path="/" element={<Home />} />
              <Route path="/select-panelists" element={<PanelistSelection />} />
              <Route path="/debate" element={<DebateGeneration />} />
              <Route path="/d/:uuid" element={<DebateViewer />} />
              <Route path="/debates" element={<DebateHistory />} />
              <Route path="*" element={<NotFound />} />
            </Routes>
          </main>
          
          <footer className="App-footer">
            <p>© 2025 Debate Generator</p>
          </footer>
        </div>
      </Router>
    </ErrorBoundary>
  );
}

export default App;
