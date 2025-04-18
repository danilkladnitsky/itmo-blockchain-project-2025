import {Card, CardProps, Text} from '@gravity-ui/uikit';
import bem from 'bem-cn-lite';

import './styles.scss';

interface SectionProps {
    children: React.ReactNode;
    theme?: CardProps['theme'];
    type?: CardProps['type'];
    view?: CardProps['view'];
    title?: string;
}

const b = bem('section');

export const Section = ({
    children,
    theme = 'normal',
    type = 'container',
    view = 'filled',
    title,
}: SectionProps) => {
    return (
        <Card theme={theme} type={type} view={view} className={b()}>
            <div className={b('header')}>
                <Text variant="body-3">{title}</Text>
            </div>
            <div className={b('content')}>{children}</div>
        </Card>
    );
};
