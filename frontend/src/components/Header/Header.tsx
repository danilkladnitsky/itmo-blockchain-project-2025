import bem from 'bem-cn-lite';
import './Header.scss';
import {Label, Link} from '@gravity-ui/uikit';

const b = bem('header');

const AppVersion = () => {
    return (
        <Link href="https://github.com/danilkladnitsky/itmo-blockchain-project-2025">
            <Label className={b('version')}>Github</Label>
        </Link>
    );
};

export const Header = () => {
    return (
        <div className={b()}>
            <div className={b('title')}>Smart Wallet</div>
            <AppVersion />
        </div>
    );
};
