import { ChangeEvent, FC, FocusEvent, ReactNode, useId } from 'react';
import styles from './TextField.module.scss';

interface TextFieldProps {
    label?: string;
    hint?: ReactNode;
    prefix?: string;
    reserveError?: boolean;
    error?: string;
    touched?: boolean;
    textarea?: boolean;
    rows?: number;
    name?: string;
    type?: string;
    placeholder?: string;
    value?: string | number;
    onChange?: (e: ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => void;
    onBlur?: (e: FocusEvent<HTMLInputElement | HTMLTextAreaElement>) => void;
    disabled?: boolean;
    required?: boolean;
    autoFocus?: boolean;
    maxLength?: number;
    minLength?: number;
    pattern?: string;
    readOnly?: boolean;
    autoComplete?: string;
    className?: string;
}

const TextField: FC<TextFieldProps> = ({
    label,
    hint,
    prefix,
    reserveError = false,
    error,
    touched,
    textarea = false,
    rows = 3,
    className = '',
    ...props
}) => {
    const id = useId();
    const hasError = Boolean(touched && error);

    const controlClassName = [
        styles.control,
        hasError ? styles.error : '',
        prefix ? styles.hasPrefix : '',
        props.disabled ? styles.disabled : '',
    ]
        .filter(Boolean)
        .join(' ');

    const fieldProps = {
        id,
        className: `${styles.input} ${textarea ? styles.textarea : ''}`,
        'aria-invalid': hasError ? true : undefined,
        ...(props as any),
    };

    return (
        <div className={`${styles.container} ${className}`}>
            {(label || hint) && (
                <div className={styles.labelRow}>
                    {label ? (
                        <label className={styles.label} htmlFor={id}>
                            {label}
                        </label>
                    ) : (
                        <span />
                    )}
                    {hint && <span className={styles.hint}>{hint}</span>}
                </div>
            )}

            <div className={controlClassName}>
                {prefix && (
                    <span className={styles.prefix} aria-hidden="true">
                        {prefix}
                    </span>
                )}
                {textarea ? <textarea rows={rows} {...fieldProps} /> : <input {...fieldProps} />}
            </div>

            {(hasError || reserveError) && (
                <span
                    className={`${styles.errorMessage} ${reserveError ? styles.reserve : ''}`}
                    role={hasError ? 'alert' : undefined}
                >
                    {hasError ? error : ''}
                </span>
            )}
        </div>
    );
};

export default TextField;
