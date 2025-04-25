import bem from 'bem-cn-lite';
import './Header.scss';
import {Label, Link} from '@gravity-ui/uikit';
import {useAppContext} from '@/App.context';

const b = bem('header');

const AppVersion = () => {
    return (
        <Link href="https://github.com/danilkladnitsky/itmo-blockchain-project-2025">
            <Label size="s" className={b('version')}>
                Github
            </Label>
        </Link>
    );
};

const BackendStatus = () => {
    const {isBackendAlive} = useAppContext();
    return (
        <Label
            size="s"
            theme={isBackendAlive ? 'success' : 'danger'}
            value={isBackendAlive ? 'live' : 'dead'}
        >
            Backend
        </Label>
    );
};

export const Header = () => {
    return (
        <div className={b()}>
            <div className={b('title')}>AI Wallet Analysis</div>
            <AppVersion />
            <BackendStatus />
        </div>
    );
};
