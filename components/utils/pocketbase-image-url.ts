const API_URL = process.env.EXPO_PUBLIC_API_URL;

export const getPocketBaseImageUrl = (
  collectionId: string,
  recordId: string,
  filename: string
) => {
  return `${API_URL}/api/files/${collectionId}/${recordId}/${filename}`;
};
