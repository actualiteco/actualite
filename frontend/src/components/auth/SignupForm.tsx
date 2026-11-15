import { useFormik } from 'formik';
import { toast } from 'react-toastify';
import TextField from '@components/shared/TextField/TextField';
import Button from '@components/shared/Button/Button';
import { useAuthStore } from '@stores/useAuthStore';
import { ApiErrorCode, parseApiError } from '@utils/error.utils';
import { isFilled, signupSchema } from './validation';
import styles from './AuthForms.module.scss';

interface SignupFormValues {
    username: string;
    email: string;
    password: string;
}

const SignupForm = () => {
    const register = useAuthStore((state) => state.register);
    const isLoading = useAuthStore((state) => state.isLoading);

    const formik = useFormik<SignupFormValues>({
        initialValues: { username: '', email: '', password: '' },
        validationSchema: signupSchema,
        onSubmit: async (values: SignupFormValues) => {
            try {
                await register(values.email.trim(), values.username, values.password);
                toast.success('Добро пожаловать в actualite!');
            } catch (error) {
                const { status, code } = parseApiError(error);

                if (status === 429) {
                    toast.error('Слишком много попыток. Попробуйте позже');
                } else if (code === ApiErrorCode.EmailAlreadyUsed) {
                    // F_Auth_6
                    formik.setFieldError('email', 'Email уже используется');
                } else if (code === ApiErrorCode.Conflict) {
                    // NF_Data_3
                    formik.setFieldError('username', 'Никнейм занят');
                } else {
                    toast.error('Не удалось зарегистрироваться. Попробуйте ещё раз');
                }
            }
        },
    });

    const shown = (name: keyof SignupFormValues) =>
        Boolean(formik.touched[name]) || formik.values[name] !== '';

    const disabled = !isFilled({ ...formik.values }) || !formik.isValid;

    return (
        <form onSubmit={formik.handleSubmit} className={styles.form} noValidate>
            <h1 className={styles.title}>Создайте аккаунт</h1>
            <p className={styles.subtitle}>Продолжайте делиться мыслями</p>

            <TextField
                name="username"
                type="text"
                placeholder="Имя пользователя"
                autoComplete="username"
                value={formik.values.username}
                onChange={formik.handleChange}
                onBlur={formik.handleBlur}
                error={formik.errors.username}
                touched={shown('username')}
            />

            <TextField
                name="email"
                type="email"
                placeholder="Email"
                autoComplete="email"
                value={formik.values.email}
                onChange={formik.handleChange}
                onBlur={formik.handleBlur}
                error={formik.errors.email}
                touched={shown('email')}
            />

            <TextField
                name="password"
                type="password"
                placeholder="Пароль"
                autoComplete="new-password"
                value={formik.values.password}
                onChange={formik.handleChange}
                onBlur={formik.handleBlur}
                error={formik.errors.password}
                touched={shown('password')}
            />

            <Button type="submit" fullWidth isLoading={isLoading} disabled={disabled}>
                Зарегистрироваться
            </Button>
        </form>
    );
};

export default SignupForm;
