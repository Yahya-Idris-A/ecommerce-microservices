package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/Yahya-idris-A/ecommerce-microservices/user-service/internal/domain"
	"github.com/google/uuid"
)

type addressUsecase struct {
	addressRepo domain.AddressRepository
}

func NewAddressUsecase(addressRepo domain.AddressRepository) domain.AddressUsecase {
	return &addressUsecase{addressRepo: addressRepo}
}

func (u *addressUsecase) AddAddress(ctx context.Context, userID uuid.UUID, req *domain.CreateAddressReq) (*domain.Address, error) {
	// 1. Cek apakah ini alamat pertama user?
	existingAddresses, err := u.addressRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Jika ini alamat pertama, paksa jadi alamat utama (Primary)
	if len(existingAddresses) == 0 {
		req.IsPrimary = true
	} else if req.IsPrimary {
		// Jika ini bukan alamat pertama tapi user mencentang "Jadikan Utama",
		// kita harus mereset semua alamat lain menjadi bukan utama dulu.
		_ = u.addressRepo.RemovePrimaryStatus(ctx, userID)
	}

	// 2. Mapping DTO ke Entitas
	address := &domain.Address{
		ID:            uuid.New(),
		UserID:        userID,
		Label:         req.Label,
		RecipientName: req.RecipientName,
		PhoneNumber:   req.PhoneNumber,
		FullAddress:   req.FullAddress,
		City:          req.City,
		Province:      req.Province,
		PostalCode:    req.PostalCode,
		IsPrimary:     req.IsPrimary,
		CreatedAt:     time.Now().UTC(),
		UpdatedAt:     time.Now().UTC(),
	}

	// 3. Simpan ke database
	err = u.addressRepo.Create(ctx, address)
	if err != nil {
		return nil, err
	}

	return address, nil
}

func (u *addressUsecase) GetMyAddresses(ctx context.Context, userID uuid.UUID) ([]domain.Address, error) {
	return u.addressRepo.GetByUserID(ctx, userID)
}

func (u *addressUsecase) UpdateAddress(ctx context.Context, userID uuid.UUID, addressID uuid.UUID, req *domain.UpdateAddressReq) (*domain.Address, error) {
	// 1. Pastikan alamatnya ada dan milik user ini
	addr, err := u.addressRepo.GetByID(ctx, addressID)
	if err != nil || addr.UserID != userID {
		return nil, errors.New("address not found or unauthorized")
	}

	// 2. Update data
	addr.Label = req.Label
	addr.RecipientName = req.RecipientName
	addr.PhoneNumber = req.PhoneNumber
	addr.FullAddress = req.FullAddress
	addr.City = req.City
	addr.Province = req.Province
	addr.PostalCode = req.PostalCode
	addr.UpdatedAt = time.Now().UTC()

	// 3. Simpan
	err = u.addressRepo.Update(ctx, addr)
	if err != nil {
		return nil, err
	}

	return addr, nil
}

func (u *addressUsecase) DeleteAddress(ctx context.Context, userID uuid.UUID, addressID uuid.UUID) error {
	// 1. Pastikan alamatnya ada dan milik user ini
	addr, err := u.addressRepo.GetByID(ctx, addressID)
	if err != nil || addr.UserID != userID {
		return errors.New("address not found or unauthorized")
	}

	// 2. Hapus alamat
	return u.addressRepo.Delete(ctx, addressID)
}

func (u *addressUsecase) SetPrimaryAddress(ctx context.Context, userID uuid.UUID, addressID uuid.UUID) error {
	// 1. Pastikan alamatnya ada dan milik user ini
	addr, err := u.addressRepo.GetByID(ctx, addressID)
	if err != nil || addr.UserID != userID {
		return errors.New("address not found or unauthorized")
	}

	// 2. Reset semua alamat user ini menjadi bukan utama
	err = u.addressRepo.RemovePrimaryStatus(ctx, userID)
	if err != nil {
		return err
	}

	// 3. Jadikan alamat yang dipilih ini sebagai alamat utama
	addr.IsPrimary = true
	addr.UpdatedAt = time.Now().UTC()

	return u.addressRepo.Update(ctx, addr)
}
