import {CurrentWallet} from '@/components/CurrentWallet/CurrentWallet';
import {WalletType} from '@/components/WalletType/WalletType';
import bem from 'bem-cn-lite';
import './styles.scss';
import {SimilarWallets} from '@/components/SimilarWallets/SimilarWallets';
import {Button} from '@gravity-ui/uikit';
import {useNavigate} from 'react-router-dom';

const b = bem('analyze-wallet-page');

export const AnalyzeWalletPage = () => {
    const navigate = useNavigate();

    return (
        <div className={b()}>
            <CurrentWallet />
            <WalletType />
            <SimilarWallets />
            <Button size="xl" view="raised" className={b('button')} onClick={() => navigate('/')}>
                Analyze another wallet
            </Button>
        </div>
    );
};
