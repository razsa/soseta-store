import PocketBase from 'pocketbase';
import { getPocketBaseImageUrl } from '../components/utils/pocketbase-image-url';
const API_URL = process.env.EXPO_PUBLIC_API_URL;

import { RecordModel } from 'pocketbase';

export type Products = RecordModel & {
  description: string;
  image: string;
  name: string;
  price: number;
  stock: string;
};
const PRODUCT_COLLECTION = 'products';
export const usePocketbase = () => {
  const pb = new PocketBase(`${API_URL}`);

  const getProducts = async () => {
    const records = await pb.collection(PRODUCT_COLLECTION).getFullList();
    return records.map((record) => ({
      ...record,
      imageUrl: getPocketBaseImageUrl(
        record.collectionId || '',
        record.id || '',
        record.image
      ),
    }));
  };

  const getProductById = async (id: string) => {
    return await pb.collection(PRODUCT_COLLECTION).getOne(id);
  }

  return {
    getProducts,
    getProductById
  };
};
