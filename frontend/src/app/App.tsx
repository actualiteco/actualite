import React from 'react';
import { ToastContainer } from 'react-toastify';
import AppRouter from '@router/AppRouter';
import ThemeToggle from '@components/shared/ThemeToggle/ThemeToggle';
import { useThemeStore } from '@stores/useThemeStore';
import './styles/fonts';
import './styles/global.scss';

const App = () => {
    const theme = useThemeStore((state) => state.theme);

    return (
        <div className="app">
            <ToastContainer
                position="bottom-right"
                hideProgressBar={true}
                theme={theme}
                autoClose={3000}
                closeOnClick
                pauseOnHover
                draggable={false}
                toastClassName="actualite-toast"
            />
            <AppRouter />
            <ThemeToggle />
        </div>
    );
};

export default App;
