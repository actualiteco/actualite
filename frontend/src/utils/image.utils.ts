export const AVATAR_MAX_BYTES = 2 * 1024 * 1024;
export const AVATAR_OUTPUT_SIZE = 400;

export interface CropArea {
    x: number;
    y: number;
    size: number;
}

export const cropToBlob = (
    image: HTMLImageElement,
    area: CropArea,
    mimeType: string,
    outputSize: number = AVATAR_OUTPUT_SIZE,
): Promise<Blob> =>
    new Promise((resolve, reject) => {
        const canvas = document.createElement('canvas');
        canvas.width = outputSize;
        canvas.height = outputSize;

        const context = canvas.getContext('2d');
        if (!context) {
            reject(new Error('Canvas is not supported'));
            return;
        }

        context.imageSmoothingQuality = 'high';
        context.drawImage(
            image,
            area.x,
            area.y,
            area.size,
            area.size,
            0,
            0,
            outputSize,
            outputSize,
        );

        canvas.toBlob(
            (blob) => (blob ? resolve(blob) : reject(new Error('Failed to create image'))),
            mimeType,
            0.92,
        );
    });
