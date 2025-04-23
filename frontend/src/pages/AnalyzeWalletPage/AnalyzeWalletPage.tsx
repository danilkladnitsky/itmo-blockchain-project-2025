import {CurrentWallet} from '@/components/CurrentWallet/CurrentWallet';
import {WalletType} from '@/components/WalletType/WalletType';
import bem from 'bem-cn-lite';
import './styles.scss';
import {SimilarWallets} from '@/components/SimilarWallets/SimilarWallets';

const b = bem('analyze-wallet-page');

export const AnalyzeWalletPage = () => {
    return (
        <div className={b()}>
            <CurrentWallet />
            <WalletType />
            <SimilarWallets />
        </div>
    );
};
