import { create } from "zustand";

export const useCart = create((set) => ({
    items: [],

    addProduct: (product: { id: string }) =>
        set((state: { items: { product: { id: string }, quantity: number }[] }) => {
            const existingItem = state.items.find((item: { product: { id: string } }) => item.product.id === product.id);
            if (existingItem) {
                return {
                    items: state.items.map((item: { product: { id: string }, quantity: number }) =>
                        item.product.id === product.id ? { ...item, quantity: item.quantity + 1 } : item
                    ),
                };
            } else {
                return {
                    items: [...state.items, { product, quantity: 1 }],
                };
            }
        }),

}));

export default useCart;
