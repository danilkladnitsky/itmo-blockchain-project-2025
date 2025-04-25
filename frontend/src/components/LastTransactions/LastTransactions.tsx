import {useAppContext} from '@/App.context';
import {Section} from '../Section/Section';
import {Box, Button, Card, Icon, Label, Link, Text} from '@gravity-ui/uikit';
import bem from 'bem-cn-lite';
import {Transaction} from '@/types/walletAnalysis';
import {
    formatAddress,
    formatDate,
    formatEther,
    formatGasPrice,
    formatHash,
} from '@/utils/formatters';

const b = bem('last-transactions');

import './styles.scss';
import {ArrowRight, Magnifier} from '@gravity-ui/icons';

const TransactionCard = ({transaction}: {transaction: Transaction}) => {
    return (
        <Card
            view="filled"
            className={b('transaction')}
            style={{
                padding: '16px',
                borderRadius: '8px',
            }}
        >
            <Box
                className={b('transaction-header')}
                style={{display: 'flex', justifyContent: 'space-between'}}
            >
                <Text variant="code-1" color="secondary">
                    {formatDate(transaction.block_timestamp)}
                </Text>
                <Label theme={transaction.receipt_status === '1' ? 'success' : 'danger'} size="s">
                    {transaction.receipt_status === '1' ? 'Success' : 'Failed'}
                </Label>
            </Box>

            <Box className={b('transaction-hash')} style={{marginBottom: '8px'}}>
                <Text variant="body-1" color="secondary">
                    Hash:{' '}
                </Text>
                <Text variant="body-1">{formatHash(transaction.hash)}</Text>
            </Box>

            <Box className={b('transaction-addresses')} style={{marginBottom: '8px'}}>
                <Box className={b('transaction-address')}>
                    <Text variant="body-1" color="secondary">
                        From:{' '}
                    </Text>
                    <Text variant="body-1">
                        {transaction.from_address_label || formatAddress(transaction.from_address)}
                    </Text>
                </Box>
                <Box className={b('transaction-address')}>
                    <Text variant="body-1" color="secondary">
                        To:{' '}
                    </Text>
                    <Text variant="body-1">
                        {transaction.to_address_label || formatAddress(transaction.to_address)}
                    </Text>
                </Box>
            </Box>

            <Box className={b('transaction-details')} style={{display: 'flex', gap: '24px'}}>
                <Box className={b('transaction-value')}>
                    <Text variant="body-1" color="secondary">
                        Value:{' '}
                    </Text>
                    <Text variant="caption-2">{formatEther(transaction.value)} ETH</Text>
                </Box>
                <Box className={b('transaction-fee')}>
                    <Text variant="body-1" color="secondary">
                        Fee:{' '}
                    </Text>
                    <Text variant="caption-2">{transaction.transaction_fee} ETH</Text>
                </Box>
                <Box className={b('transaction-gas')}>
                    <Text variant="body-1" color="secondary">
                        Gas Price:{' '}
                    </Text>
                    <Text variant="caption-2">{formatGasPrice(transaction.gas_price)} Gwei</Text>
                </Box>
            </Box>
            <Box
                style={{
                    display: 'flex',
                    flexDirection: 'column',
                    gap: '8px',
                    marginTop: '16px',
                    justifyContent: 'flex-end',
                }}
            >
                <Link href={`/?wallet=${transaction.from_address}`}>
                    <Button size="s" view="action" width="max">
                        Analyze from address
                        <Icon data={Magnifier} />
                    </Button>
                </Link>
                <Link href={`/?wallet=${transaction.to_address}`}>
                    <Button size="s" view="action" width="max">
                        Analyze to address
                        <Icon data={Magnifier} />
                    </Button>
                </Link>
                <Link href={`https://etherscan.io/tx/${transaction.hash}`} target="_blank">
                    <Button size="s" view="normal" width="max">
                        View on explorer
                        <Icon data={ArrowRight} />
                    </Button>
                </Link>
            </Box>
        </Card>
    );
};

export const LastTransactions = () => {
    const {walletAnalysis} = useAppContext();
    const {transactions} = walletAnalysis || {};

    return (
        <Section title="Wallet transactions">
            <Box className={b()}>
                <Text variant="body-1">
                    Last transactions: <Text color="secondary">{transactions?.length || 0}</Text>
                </Text>
                <Box className={b('transactions')}>
                    {transactions?.map((transaction) => (
                        <TransactionCard key={transaction.hash} transaction={transaction} />
                    ))}
                </Box>
            </Box>
        </Section>
    );
};
