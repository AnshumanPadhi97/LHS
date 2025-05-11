import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import Layout from './components/layout/Layout';

function App() {
    return (
       <Router>
        <Routes>
        <Route path="/" element={<Layout />}>
            {/* <Route index element={<Dashboard />} />
            <Route path="new-stack" element={<NewStack />} />
            <Route path="customize-stack/:templateId" element={<CustomizeStack />} />
            <Route path="stacks/:stackId" element={<StackDetails />} />
            <Route path="logs/:stackId/:serviceId" element={<LogsPage />} />
            <Route path="settings" element={<Settings />} /> */}
        </Route>
        </Routes>
    </Router>
    )
}

export default App
