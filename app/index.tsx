import { Text, View, FlatList } from "react-native";
import products from "../assets/products.json";
import ProductListItem from "../components/ProductListItem";

export default function Index() {
  return (
    <View
      style={{
        flex: 1,
        justifyContent: "center",
        alignItems: "center",
      }}
    >
      <FlatList
        data={products}
        renderItem={({ item }) => (
          <ProductListItem product={item} />
        )}
      />
    </View>
  );
}
