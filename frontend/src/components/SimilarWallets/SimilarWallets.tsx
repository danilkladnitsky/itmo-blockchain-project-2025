import {Section} from '../Section/Section';
import {Box, Card, ClipboardButton, Progress, ProgressTheme, Text} from '@gravity-ui/uikit';
import bem from 'bem-cn-lite';
import {getWalletAddress} from '@/utils/formatWalletAddress';

const b = bem('similar-wallets');

import './styles.scss';
interface WalletProps {
    address: string;
    balance: number;
}
const Wallet = ({address, balance}: WalletProps) => {
    const confidence = Math.round(Math.random() * 100);
    const confidenceTheme: ProgressTheme =
        confidence > 80 ? 'success' : confidence > 60 ? 'warning' : 'default';

    const copyToClipboard = () => {
        navigator.clipboard.writeText(address);
    };

    return (
        <Card className={b('wallet')}>
            <Text variant="body-1">
                Address:{' '}
                <Box>
                    <Text color="secondary">{getWalletAddress(address)}</Text>
                    <ClipboardButton
                        onClick={copyToClipboard}
                        content={address}
                        text="Copy"
                        size="xs"
                    />
                </Box>
            </Text>
            <Text variant="body-1">
                Balance: <Text color="secondary">{balance}</Text>
            </Text>
            <Box>
                <Progress
                    text={`${confidence}% confidence`}
                    value={confidence}
                    theme={confidenceTheme}
                />
            </Box>
        </Card>
    );
};

export const SimilarWallets = () => {
    return (
        <Section title="Similar wallets" view="filled">
            <Box className={b()}>
                <Wallet address="0x1234567890123456789012345678901234567890" balance={100} />
                <Wallet address="0x1234567890123456789012345678901234567890" balance={100} />
                <Wallet address="0x1234567890123456789012345678901234567890" balance={100} />
            </Box>
        </Section>
    );
};
