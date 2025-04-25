/* eslint-disable camelcase */
import {Section} from '../Section/Section';
import {Box, Progress, Text} from '@gravity-ui/uikit';
import {useAppContext} from '@/App.context';
import bem from 'bem-cn-lite';

const b = bem('wallet-type');

import './styles.scss';

export const WalletType = () => {
    const {walletAnalysis} = useAppContext();

    const {predicted_class, confidence} = walletAnalysis?.ml_result || {};
    const confidenceValue = confidence ? Math.round(confidence * 100) : 0;

    return (
        <Section title="Wallet type">
            <Box className={b()}>
                <Text color="secondary" variant="subheader-1">
                    Wallet type:{' '}
                </Text>
                <Text variant="body-1">{predicted_class || 'Unknown'}</Text>
                <Box className={b('progress')}>
                    <Progress
                        size="m"
                        text={`${confidenceValue}% confidence`}
                        value={confidenceValue}
                        theme="warning"
                    />
                </Box>
            </Box>
        </Section>
    );
};
