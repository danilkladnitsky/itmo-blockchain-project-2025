import {Box, ClipboardButton, Text} from '@gravity-ui/uikit';
import {useAppContext} from '../../App.context';
import bem from 'bem-cn-lite';

import {Section} from '../Section/Section';
import {getWalletAddress} from '@/utils/formatWalletAddress';
const b = bem('current-wallet');

import './styles.scss';

export const CurrentWallet = () => {
    const {walletAddress} = useAppContext();

    const copyToClipboard = () => {
        navigator.clipboard.writeText(walletAddress);
    };

    return (
        <Section title="Wallet address">
            <Text ellipsis variant="body-3" className={b()}>
                <Text color="secondary" variant="subheader-1">
                    Wallet address:{' '}
                </Text>
                <Box className={b('address')}>
                    <Text variant="body-1">{getWalletAddress(walletAddress)}</Text>
                    <ClipboardButton
                        onClick={copyToClipboard}
                        content={walletAddress}
                        text="Copy"
                        size="s"
                    />
                </Box>
            </Text>
        </Section>
    );
};
