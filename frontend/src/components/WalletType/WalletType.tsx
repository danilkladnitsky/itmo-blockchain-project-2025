import {Section} from '../Section/Section';
import {Text} from '@gravity-ui/uikit';
import {useAppContext} from '@/App.context';

export const WalletType = () => {
    const {walletType} = useAppContext();

    return (
        <Section title="Wallet type">
            <Text color="secondary" variant="subheader-1">
                Wallet type:{' '}
            </Text>
            <Text variant="body-1">{walletType || 'Unknown'}</Text>
        </Section>
    );
};
