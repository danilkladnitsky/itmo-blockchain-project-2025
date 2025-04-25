import {Box, Button, Icon, Text, TextInput} from '@gravity-ui/uikit';
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

const tryWallets = [
    '0xcB1C1FdE09f811B294172696404e88E658659905',
    '0x1231deb6f5749ef6ce6943a275a1d3e7486f4eae',
];

export const InsertWalletPage = () => {
    const {walletAddress, setWalletAddress, searchWallet, isLoadingWalletAnalysis, error} =
        useAppContext();

    const location = useLocation();

    const handleSearch = () => {
        searchWallet(walletAddress);
    };

    const handleLastWalletClick = (wallet: string) => {
        setWalletAddress(wallet);
        searchWallet(wallet);
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
                        value={walletAddress}
                        defaultValue={extractWalletAddressFromParams(location.search)}
                        onChange={(e) => setWalletAddress(e.target.value)}
                        disabled={isLoadingWalletAnalysis}
                        error={error}
                        errorPlacement="inside"
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
                <Box>
                    <Text variant="body-3" className={b('form-title')}>
                        Try these wallets
                    </Text>
                    <div className={b('try-wallets-list')}>
                        {tryWallets.map((wallet) => (
                            <Button
                                key={wallet}
                                view="raised"
                                className={b('try-wallets-item')}
                                onClick={() => handleLastWalletClick(wallet)}
                            >
                                {wallet}
                            </Button>
                        ))}
                    </div>
                </Box>
            </div>
        </div>
    );
};
