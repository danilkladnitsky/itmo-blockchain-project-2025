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
    const {
        walletAddress,
        setWalletAddress,
        searchWallet,
        isLoadingWalletAnalysis,
        lastWalletsList,
    } = useAppContext();

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
                <div
                    style={{display: lastWalletsList.length > 0 ? 'block' : 'none'}}
                    className={b('last-wallets')}
                >
                    <Text variant="body-3" className={b('last-wallets-title')}>
                        Last wallets
                    </Text>
                    <div
                        className={b('last-wallets-list')}
                        style={{
                            display: 'flex',
                            flexDirection: 'column',
                            marginTop: '16px',
                            gap: '8px',
                        }}
                    >
                        {lastWalletsList.map((wallet) => (
                            <Button
                                key={wallet}
                                width="max"
                                view="raised"
                                className={b('last-wallets-item')}
                                onClick={() => handleLastWalletClick(wallet)}
                            >
                                {wallet}
                            </Button>
                        ))}
                    </div>
                </div>
            </div>
        </div>
    );
};
