package usecase

import (
	"context"
	"time"

	"github.com/Yahya-idris-A/ecommerce-microservices/order-service/internal/domain"
	"github.com/google/uuid"
)

// ==========================================
// KONTRAK UNTUK KOMUNIKASI ANTAR SERVICE
// ==========================================
// Nanti kita buat implementasinya untuk nembak API Product Service
type ProductServiceClient interface {
	GetProductDetails(ctx context.Context, productID uuid.UUID) (name string, img string, price float64, err error)
	GetMerchantName(ctx context.Context, merchantID uuid.UUID) (name string, err error)
}

type cartUsecase struct {
	cartRepo      domain.CartRepository
	productClient ProductServiceClient // Inject service eksternal
}

func NewCartUsecase(cartRepo domain.CartRepository, productClient ProductServiceClient) domain.CartUsecase {
	return &cartUsecase{
		cartRepo:      cartRepo,
		productClient: productClient,
	}
}

// ==========================================
// IMPLEMENTASI LOGIKA BISNIS
// ==========================================

func (u *cartUsecase) AddItemToCart(ctx context.Context, userID uuid.UUID, req *domain.AddCartItemReq) error {
	// 1. Cek apakah user sudah punya payung keranjang?
	cart, err := u.cartRepo.GetCartByUserID(ctx, userID)
	if err != nil {
		if err.Error() == "cart not found" {
			// Auto-Create Cart kalau belum punya
			cart = &domain.Cart{
				ID:        uuid.New(),
				UserID:    userID,
				CreatedAt: time.Now().UTC(),
				UpdatedAt: time.Now().UTC(),
			}
			if err := u.cartRepo.CreateCart(ctx, cart); err != nil {
				return err
			}
		} else {
			return err
		}
	}

	productID, _ := uuid.Parse(req.ProductID)
	merchantID, _ := uuid.Parse(req.MerchantID)

	// 2. Cek apakah barang ini sudah ada di keranjang?
	existingItem, err := u.cartRepo.GetItemByCartAndProduct(ctx, cart.ID, productID)
	if err != nil {
		return err
	}

	if existingItem != nil {
		// Logika Upsert: Kalau barang sudah ada, cukup tambahkan quantity-nya
		newQuantity := existingItem.Quantity + req.Quantity
		return u.cartRepo.UpdateItemQuantity(ctx, existingItem.ID, newQuantity)
	}

	// 3. Kalau barang belum ada, buat item baru
	newItem := &domain.CartItem{
		ID:         uuid.New(),
		CartID:     cart.ID,
		ProductID:  productID,
		MerchantID: merchantID,
		Quantity:   req.Quantity,
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	}

	return u.cartRepo.AddItem(ctx, newItem)
}

func (u *cartUsecase) GetMyCart(ctx context.Context, userID uuid.UUID) (*domain.CartResponse, error) {
	// 1. Cari keranjangnya
	cart, err := u.cartRepo.GetCartByUserID(ctx, userID)
	if err != nil {
		if err.Error() == "cart not found" {
			// Kalau belum punya keranjang, kembalikan response kosong yang rapi
			return &domain.CartResponse{
				UserID:      userID,
				Groups:      []domain.MerchantCartGroup{},
				TotalAmount: 0,
			}, nil
		}
		return nil, err
	}

	// 2. Ambil semua barang di keranjang
	items, err := u.cartRepo.GetItemsByCartID(ctx, cart.ID)
	if err != nil {
		return nil, err
	}

	// 3. Logika Grouping per Toko (Merchant)
	groupsMap := make(map[uuid.UUID]*domain.MerchantCartGroup)
	var grandTotal float64

	for _, item := range items {
		// PANGGIL MICROSERVICE SEBELAH (Product Service)
		// Catatan: Di realita unicorn, pemanggilan ini biasanya di-batch/digabung
		// agar tidak nembak API berkali-kali di dalam loop. Tapi untuk portofolio, ini sudah mantap.
		prodName, prodImg, prodPrice, _ := u.productClient.GetProductDetails(ctx, item.ProductID)

		subTotal := prodPrice * float64(item.Quantity)
		grandTotal += subTotal

		itemDetail := domain.CartItemDetail{
			ItemID:      item.ID,
			ProductID:   item.ProductID,
			ProductName: prodName,
			ProductImg:  prodImg,
			Price:       prodPrice,
			Quantity:    item.Quantity,
			SubTotal:    subTotal,
		}

		// Masukkan ke dalam grup tokonya
		if group, exists := groupsMap[item.MerchantID]; exists {
			group.Items = append(group.Items, itemDetail)
		} else {
			merchName, _ := u.productClient.GetMerchantName(ctx, item.MerchantID)
			groupsMap[item.MerchantID] = &domain.MerchantCartGroup{
				MerchantID:   item.MerchantID,
				MerchantName: merchName,
				Items:        []domain.CartItemDetail{itemDetail},
			}
		}
	}

	// 4. Ubah Map menjadi Array/Slice agar format JSON-nya cakep
	var groups []domain.MerchantCartGroup
	for _, group := range groupsMap {
		groups = append(groups, *group)
	}

	return &domain.CartResponse{
		CartID:      cart.ID,
		UserID:      userID,
		Groups:      groups,
		TotalAmount: grandTotal,
	}, nil
}

func (u *cartUsecase) UpdateItemQuantity(ctx context.Context, userID uuid.UUID, itemID uuid.UUID, req *domain.UpdateCartItemReq) error {
	// Di sistem nyata, kita bisa tambahkan validasi apakah itemID ini
	// benar-benar ada di dalam cart milik userID.
	return u.cartRepo.UpdateItemQuantity(ctx, itemID, req.Quantity)
}

func (u *cartUsecase) RemoveItemFromCart(ctx context.Context, userID uuid.UUID, itemID uuid.UUID) error {
	return u.cartRepo.DeleteItem(ctx, itemID)
}
