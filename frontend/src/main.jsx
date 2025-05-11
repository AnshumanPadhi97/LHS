import {createRoot} from 'react-dom/client'
import App from './App'
import './App.css';
import {HeroUIProvider} from "@heroui/react";

createRoot(document.getElementById('root')).render(
    <HeroUIProvider>
        <App/>
    </HeroUIProvider>
)
