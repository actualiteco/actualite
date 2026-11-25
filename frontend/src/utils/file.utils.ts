export const MB = 1024 * 1024;

const IMAGE_EXTENSION = /\.(jpe?g|png)$/i;
const IMAGE_MIME_TYPES = ['image/jpeg', 'image/png'];

export const validateImageFile = (file: File, maxBytes: number): string | null => {
    if (!IMAGE_EXTENSION.test(file.name) || !IMAGE_MIME_TYPES.includes(file.type)) {
        return 'Допустимы только изображения JPG или PNG';
    }
    if (file.size > maxBytes) {
        return `Размер файла не должен превышать ${Math.round(maxBytes / MB)} МБ`;
    }
    return null;
};
