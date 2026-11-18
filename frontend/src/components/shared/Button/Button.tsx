import React, { FC, ButtonHTMLAttributes } from 'react';
import styles from './Button.module.scss';

interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
    variant?: 'contained' | 'outlined' | 'text';
    fullWidth?: boolean;
    isLoading?: boolean;
}

const Button: FC<ButtonProps> = ({
    children,
    variant = 'contained',
    fullWidth = false,
    isLoading = false,
    disabled,
    type = 'button',
    ...props
}) => {
    return (
        <button
            type={type}
            className={`${styles.button} ${styles[variant]} ${fullWidth ? styles.fullWidth : ''}`}
            disabled={disabled || isLoading}
            {...props}
        >
            {isLoading ? 'Загрузка…' : children}
        </button>
    );
};

export default Button;
