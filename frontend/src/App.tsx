import {ThemeProvider} from '@gravity-ui/uikit';
import bem from 'bem-cn-lite';

import {DEFAULT_THEME} from './constants';
import {InsertWalletPage} from './pages/InsertWalletPage/InsertWalletPage';

import './main.scss';
import {AppContextProvider} from './App.context';
import {Header} from './components/Header/Header';

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
        <AppWrapper>
            <Header />
            <InsertWalletPage />
        </AppWrapper>
    );
};

export default App;
