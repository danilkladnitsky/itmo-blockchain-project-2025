import {Section} from '../Section/Section';
import {Box, Progress, Text} from '@gravity-ui/uikit';
import {useAppContext} from '@/App.context';
import bem from 'bem-cn-lite';

const b = bem('wallet-type');

import './styles.scss';

export const WalletType = () => {
    const {walletType} = useAppContext();

    return (
        <Section title="Wallet type">
            <Box className={b()}>
                <Text color="secondary" variant="subheader-1">
                    Wallet type:{' '}
                </Text>
                <Text variant="body-1">{walletType || 'Unknown'}</Text>
                <Box className={b('progress')}>
                    <Progress size="m" text="85% confidence" value={85} theme="warning" />
                </Box>
            </Box>
        </Section>
    );
};
