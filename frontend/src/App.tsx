import {ThemeProvider} from '@gravity-ui/uikit';
import bem from 'bem-cn-lite';
import {BrowserRouter, Route, Routes} from 'react-router-dom';

import {DEFAULT_THEME} from './constants';
import {InsertWalletPage} from './pages/InsertWalletPage/InsertWalletPage';

import './main.scss';
import {AppContextProvider} from './App.context';
import {Header} from './components/Header/Header';
import {AnalyzeWalletPage} from './pages/AnalyzeWalletPage/AnalyzeWalletPage';

interface AppWrapperProps {
    children: React.ReactNode;
}

const b = bem('app');

const AppWrapper = ({children}: AppWrapperProps) => {
    return (
        <AppContextProvider>
            <ThemeProvider theme={DEFAULT_THEME}>
                <div className={b()}>
                    <div className={b('wrapper')}>{children}</div>
                </div>
            </ThemeProvider>
        </AppContextProvider>
    );
};

const App = () => {
    return (
        <BrowserRouter>
            <AppWrapper>
                <Header />
                <Routes>
                    <Route path="/" element={<InsertWalletPage />} />
                    <Route path="/analyze" element={<AnalyzeWalletPage />} />
                </Routes>
            </AppWrapper>
        </BrowserRouter>
    );
};

export default App;
