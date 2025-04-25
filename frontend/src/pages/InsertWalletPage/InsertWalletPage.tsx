import {Button, Icon, Text, TextInput} from '@gravity-ui/uikit';
import bem from 'bem-cn-lite';

import './styles.scss';
import {CreditCard} from '@gravity-ui/icons';
import {useAppContext} from '@/App.context';
import {useLocation} from 'react-router';
import {useEffect} from 'react';

const b = bem('insert-wallet-page');
const extractWalletAddressFromParams = (searchParams: string) => {
    const params = new URLSearchParams(searchParams);
    return params.get('wallet') || '';
};

export const InsertWalletPage = () => {
    const {walletAddress, setWalletAddress, searchWallet, isLoadingWalletAnalysis} =
        useAppContext();

    const location = useLocation();

    const handleSearch = () => {
        searchWallet();
    };

    useEffect(() => {
        setWalletAddress(extractWalletAddressFromParams(location.search));
    }, [location.search]);

    return (
        <div className={b()}>
            <div className={b('form')}>
                <Text variant="body-3" className={b('form-title')}>
                    Provide existing wallet address <Icon size={20} data={CreditCard} />
                </Text>
                <div className={b('form-input')}>
                    <TextInput
                        placeholder="Wallet address in hex format"
                        size="xl"
                        hasClear
                        value={walletAddress}
                        defaultValue={extractWalletAddressFromParams(location.search)}
                        onChange={(e) => setWalletAddress(e.target.value)}
                        disabled={isLoadingWalletAnalysis}
                    />
                    <Button
                        loading={isLoadingWalletAnalysis}
                        size="xl"
                        view="action"
                        onClick={handleSearch}
                        disabled={isLoadingWalletAnalysis}
                    >
                        Search
                    </Button>
                </div>
            </div>
        </div>
    );
};
