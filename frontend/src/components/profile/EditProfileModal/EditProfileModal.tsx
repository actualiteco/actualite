import { ChangeEvent, FC, ReactNode, useState } from 'react';
import { useFormik } from 'formik';
import * as Yup from 'yup';
import { toast } from 'react-toastify';
import Modal from '@components/shared/Modal/Modal';
import TextField from '@components/shared/TextField/TextField';
import Button from '@components/shared/Button/Button';
import type { Profile, ProfileChanges } from '@models/Profile';
import { getDisplayName } from '@models/Profile';
import { useProfileStore } from '@stores/useProfileStore';
import { ApiErrorCode, parseApiError } from '@utils/error.utils';
import {
    DISPLAY_NAME_MAX,
    STATUS_MAX,
    USERNAME_MAX,
    charCount,
    displayNameField,
    isBlank,
    sanitizeUsername,
    statusField,
    truncateChars,
    usernameField,
} from '@utils/validators.utils';
import styles from './EditProfileModal.module.scss';

interface FormValues {
    displayName: string;
    username: string;
    status: string;
}

interface EditProfileModalProps {
    open: boolean;
    profile: Profile;
    onClose: () => void;
}

const buildChanges = (values: FormValues, initial: FormValues): ProfileChanges => {
    const changes: ProfileChanges = {};
    const displayName = values.displayName.trim();
    const status = values.status.trim();

    if (displayName !== initial.displayName.trim()) {
        changes.displayName = displayName;
    }
    if (values.username !== initial.username) {
        changes.username = values.username;
    }
    if (status !== initial.status.trim()) {
        changes.status = status;
    }
    return changes;
};

const buildSchema = (initial: FormValues) =>
    Yup.object().shape({
        displayName: displayNameField(),
        username: usernameField(),
        status: statusField().test(
            'not-cleared',
            'Статус не может быть пустым',
            (value) => isBlank(initial.status) || !isBlank(value),
        ),
    });

interface CountedProps {
    value: string;
    max: number;
    children: ReactNode;
}

const Counted: FC<CountedProps> = ({ value, max, children }) => (
    <div className={styles.field}>
        {children}
        <span className={styles.counter}>
            {charCount(value)}/{max}
        </span>
    </div>
);

const EditProfileForm: FC<Omit<EditProfileModalProps, 'open'>> = ({ profile, onClose }) => {
    const updateProfile = useProfileStore((state) => state.updateProfile);

    const [initialValues] = useState<FormValues>(() => ({
        displayName: getDisplayName(profile),
        username: profile.username,
        status: profile.status,
    }));
    const [schema] = useState(() => buildSchema(initialValues));

    const formik = useFormik<FormValues>({
        initialValues,
        validationSchema: schema,
        onSubmit: async (values, { setFieldError }) => {
            const changes = buildChanges(values, initialValues);
            if (Object.keys(changes).length === 0) {
                onClose();
                return;
            }

            try {
                await updateProfile(changes);
                toast.success('Профиль обновлён');
                onClose();
            } catch (error) {
                const { status, code, isNetwork } = parseApiError(error);
                if (isNetwork) {
                    return;
                }
                if (code === ApiErrorCode.Conflict) {
                    setFieldError('username', 'Никнейм занят');
                } else if (status === 429) {
                    toast.error('Слишком много попыток. Попробуйте позже');
                } else {
                    toast.error('Не удалось сохранить профиль. Попробуйте ещё раз');
                }
            }
        },
    });

    const handleChange =
        (name: keyof FormValues, transform: (value: string) => string) =>
        (event: ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
            formik.setFieldValue(name, transform(event.target.value));
        };

    const shown = (name: keyof FormValues) =>
        Boolean(formik.touched[name]) || formik.values[name] !== initialValues[name];

    const hasChanges = Object.keys(buildChanges(formik.values, initialValues)).length > 0;
    const disabled = !hasChanges || !formik.isValid || formik.isSubmitting;

    return (
        <form onSubmit={formik.handleSubmit} className={styles.form} noValidate>
            <h2 className={styles.title}>Редактировать профиль</h2>

            <Counted value={formik.values.displayName} max={DISPLAY_NAME_MAX}>
                <TextField
                    label="Имя"
                    name="displayName"
                    autoComplete="off"
                    value={formik.values.displayName}
                    onChange={handleChange('displayName', (value) =>
                        truncateChars(value, DISPLAY_NAME_MAX),
                    )}
                    onBlur={formik.handleBlur}
                    error={formik.errors.displayName}
                    touched={shown('displayName')}
                />
            </Counted>

            <Counted value={formik.values.username} max={USERNAME_MAX}>
                <TextField
                    label="Никнейм"
                    name="username"
                    autoComplete="off"
                    value={formik.values.username}
                    onChange={handleChange('username', sanitizeUsername)}
                    onBlur={formik.handleBlur}
                    error={formik.errors.username}
                    touched={shown('username')}
                />
            </Counted>

            <Counted value={formik.values.status} max={STATUS_MAX}>
                <TextField
                    label="Статус"
                    name="status"
                    autoComplete="off"
                    value={formik.values.status}
                    onChange={handleChange('status', (value) => truncateChars(value, STATUS_MAX))}
                    onBlur={formik.handleBlur}
                    error={formik.errors.status}
                    touched={shown('status')}
                />
            </Counted>

            <div className={styles.actions}>
                <Button variant="outlined" onClick={onClose}>
                    Отмена
                </Button>
                <Button type="submit" isLoading={formik.isSubmitting} disabled={disabled}>
                    Сохранить
                </Button>
            </div>
        </form>
    );
};

const EditProfileModal: FC<EditProfileModalProps> = ({ open, profile, onClose }) => (
    <Modal open={open} onClose={onClose}>
        <EditProfileForm profile={profile} onClose={onClose} />
    </Modal>
);

export default EditProfileModal;
