import {CurrentWallet} from '@/components/CurrentWallet/CurrentWallet';
import {WalletType} from '@/components/WalletType/WalletType';
import bem from 'bem-cn-lite';
import './styles.scss';

const b = bem('analyze-wallet-page');

export const AnalyzeWalletPage = () => {
    return (
        <div className={b()}>
            <CurrentWallet />
            <WalletType />
        </div>
    );
};
