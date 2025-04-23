import {Section} from '../Section/Section';
import {Box, Card, Progress, ProgressTheme, Text} from '@gravity-ui/uikit';
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

    return (
        <Card className={b('wallet')}>
            <Text variant="body-1">Address: {getWalletAddress(address)}</Text>
            <Text variant="body-1">Balance: {balance}</Text>
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
