import {Text} from '@gravity-ui/uikit';
import {useAppContext} from '../../App.context';
import bem from 'bem-cn-lite';

import './styles.scss';
import {Section} from '../Section/Section';
import {getWalletAddress} from '@/utils/formatWalletAddress';

const b = bem('current-wallet');

export const CurrentWallet = () => {
    const {walletAddress} = useAppContext();

    return (
        <Section title="Wallet address">
            <Text ellipsis variant="body-3" className={b()}>
                <Text color="secondary" variant="subheader-1">
                    Wallet address:{' '}
                </Text>
                <Text variant="body-1">{getWalletAddress(walletAddress)}</Text>
            </Text>
        </Section>
    );
};
