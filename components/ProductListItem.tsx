import { Text, View } from "react-native";
export default function ProductListItem({ product }: { product: Product }) {
  return (
    <View>
      <Text>{product.name}</Text>
    </View>
  );
}
