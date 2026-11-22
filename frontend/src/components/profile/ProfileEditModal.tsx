import React, { FC, useEffect, useRef, useState } from 'react';
import { useFormik } from 'formik';
import { toast } from 'react-toastify';
import Avatar from '@components/shared/Avatar/Avatar';
import Button from '@components/shared/Button/Button';
import Modal from '@components/shared/Modal/Modal';
import TextField from '@components/shared/TextField/TextField';
import { ProfileService } from '@services/ProfileService';
import type { Profile, ProfileChanges } from '@models/Profile';
import { getDisplayName } from '@models/Profile';
import { getApiErrorMessage, parseApiError } from '@utils/error.utils';
import { validateImageFile } from '@utils/file.utils';
import { charCount } from '@utils/validators.utils';
import AvatarCropModal, { AVATAR_MAX_BYTES } from './AvatarCropModal';
import { STATUS_MAX, isUsernameFormatValid, profileSchema } from './validation';
import styles from './ProfileEditModal.module.scss';

interface ProfileEditModalProps {
    open: boolean;
    profile: Profile;
    onClose: () => void;
    onSaved: (profile: Profile) => void;
}

interface ProfileFormValues {
    displayName: string;
    username: string;
    status: string;
}

const toValues = (profile: Profile): ProfileFormValues => ({
    displayName: getDisplayName(profile),
    username: profile.username,
    status: profile.status,
});

const buildChanges = (values: ProfileFormValues, initial: ProfileFormValues): ProfileChanges => {
    const changes: ProfileChanges = {};
    const displayName = values.displayName.trim();
    if (displayName !== initial.displayName.trim()) {
        changes.displayName = displayName;
    }
    if (values.username !== initial.username) {
        changes.username = values.username;
    }
    if (values.status.trim() !== initial.status.trim()) {
        changes.status = values.status.trim();
    }
    return changes;
};

const ProfileEditModal: FC<ProfileEditModalProps> = ({ open, profile, onClose, onSaved }) => {
    const [cropFile, setCropFile] = useState<File | null>(null);
    const [avatarError, setAvatarError] = useState<string | null>(null);
    const [isAvatarSaving, setIsAvatarSaving] = useState(false);
    const fileInput = useRef<HTMLInputElement>(null);

    const initial = toValues(profile);

    const formik = useFormik<ProfileFormValues>({
        initialValues: initial,
        validationSchema: profileSchema,
        onSubmit: async (values: ProfileFormValues) => {
            const changes = buildChanges(values, initial);

            if (Object.keys(changes).length === 0) {
                onClose();
                return;
            }
            try {
                const updated = await ProfileService.update(changes);
                toast.success('Профиль обновлён');
                onSaved({
                    ...profile,
                    displayName:
                        changes.displayName !== undefined
                            ? updated.displayName || changes.displayName
                            : profile.displayName,
                    username:
                        changes.username !== undefined
                            ? updated.username || changes.username
                            : profile.username,
                    status: changes.status !== undefined ? changes.status : profile.status,
                });
                onClose();
            } catch (error) {
                if (parseApiError(error).status === 409) {
                    formik.setFieldError('username', 'Никнейм занят');
                } else {
                    toast.error(getApiErrorMessage(error, 'Не удалось сохранить профиль'));
                }
            }
        },
    });

    useEffect(() => {
        if (open) {
            formik.resetForm({ values: toValues(profile) });
            setAvatarError(null);
        }
    }, [open, profile]);

    const handleUsernameBlur = async (event: React.FocusEvent<HTMLInputElement | HTMLTextAreaElement>) => {
        formik.handleBlur(event);
        const username = formik.values.username;
        if (username === profile.username || !isUsernameFormatValid(username)) {
            return;
        }
        try {
            if (!(await ProfileService.isUsernameAvailable(username))) {
                formik.setFieldError('username', 'Никнейм занят');
            }
        } catch {
            // Проверка необязательна: окончательно занятость проверит сервер при сохранении
        }
    };

    const handleFile = (event: React.ChangeEvent<HTMLInputElement>) => {
        const file = event.target.files?.[0];
        event.target.value = '';
        if (!file) {
            return;
        }
        const error = validateImageFile(file, AVATAR_MAX_BYTES);
        setAvatarError(error);
        if (!error) {
            setCropFile(file);
        }
    };

    const handleAvatarConfirm = async (blob: Blob, fileName: string) => {
        setIsAvatarSaving(true);
        try {
            const avatarUrl = await ProfileService.uploadAvatar(blob, fileName);
            toast.success('Аватар обновлён');
            onSaved({ ...profile, avatarUrl });
            setCropFile(null);
        } catch (error) {
            toast.error(getApiErrorMessage(error, 'Не удалось загрузить аватар'));
        } finally {
            setIsAvatarSaving(false);
        }
    };

    const shown = (name: keyof ProfileFormValues) =>
        Boolean(formik.touched[name]) || formik.values[name] !== initial[name];
    const hasChanges = Object.keys(buildChanges(formik.values, initial)).length > 0;
    const disabled = !hasChanges || !formik.isValid || formik.isSubmitting;

    return (
        <>
            <Modal open={open && !cropFile} onClose={onClose}>
                <form className={styles.form} onSubmit={formik.handleSubmit} noValidate>
                    <h1 className={styles.title}>Редактировать профиль</h1>

                    <div className={styles.avatarRow}>
                        <Avatar url={profile.avatarUrl} name={getDisplayName(profile)} size={64} />
                        <div>
                            <input
                                ref={fileInput}
                                className={styles.hiddenInput}
                                type="file"
                                accept=".jpg,.jpeg,.png,image/jpeg,image/png"
                                onChange={handleFile}
                            />
                            <Button variant="outlined" onClick={() => fileInput.current?.click()}>
                                Сменить фото
                            </Button>
                            <p className={styles.avatarHint}>JPG или PNG, до 2 МБ</p>
                            {avatarError && <p className={styles.avatarError}>{avatarError}</p>}
                        </div>
                    </div>

                    <TextField
                        name="displayName"
                        placeholder="Имя"
                        autoComplete="name"
                        value={formik.values.displayName}
                        onChange={formik.handleChange}
                        onBlur={formik.handleBlur}
                        error={formik.errors.displayName}
                        touched={shown('displayName')}
                    />
                    <TextField
                        name="username"
                        placeholder="Никнейм"
                        autoComplete="username"
                        value={formik.values.username}
                        onChange={formik.handleChange}
                        onBlur={handleUsernameBlur}
                        error={formik.errors.username}
                        touched={shown('username')}
                    />
                    <TextField
                        name="status"
                        placeholder="Статус"
                        textarea
                        rows={3}
                        value={formik.values.status}
                        onChange={formik.handleChange}
                        onBlur={formik.handleBlur}
                        error={formik.errors.status}
                        touched={shown('status')}
                    />
                    <span className={styles.counter}>
                        {charCount(formik.values.status)}/{STATUS_MAX}
                    </span>

                    <div className={styles.actions}>
                        <Button variant="outlined" onClick={onClose}>
                            Отмена
                        </Button>
                        <Button type="submit" disabled={disabled} isLoading={formik.isSubmitting}>
                            Сохранить
                        </Button>
                    </div>
                </form>
            </Modal>

            <AvatarCropModal
                file={cropFile}
                isSaving={isAvatarSaving}
                onCancel={() => setCropFile(null)}
                onConfirm={handleAvatarConfirm}
            />
        </>
    );
};

export default ProfileEditModal;
