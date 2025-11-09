package repository

import (
	"context"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TestUser is a test entity
type TestUser struct {
	ID    uint   `gorm:"primarykey"`
	Name  string `gorm:"size:100"`
	Email string `gorm:"size:100;uniqueIndex"`
	Age   int
}

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	// Use SQLite in-memory database for testing
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Auto-migrate the test schema
	if err := db.AutoMigrate(&TestUser{}); err != nil {
		t.Fatalf("Failed to migrate test schema: %v", err)
	}

	return db
}

func TestNew(t *testing.T) {
	db := setupTestDB(t)
	repo := New[TestUser](db)

	if repo == nil {
		t.Fatal("Expected repository to be created")
	}
	if repo.db == nil {
		t.Fatal("Expected repository to have database connection")
	}
}

func TestCreate(t *testing.T) {
	db := setupTestDB(t)
	repo := New[TestUser](db)
	ctx := context.Background()

	t.Run("creates new record successfully", func(t *testing.T) {
		user := &TestUser{
			Name:  "John Doe",
			Email: "john@example.com",
			Age:   30,
		}

		err := repo.Create(ctx, user)
		if err != nil {
			t.Fatalf("Failed to create user: %v", err)
		}
		if user.ID == 0 {
			t.Error("Expected user ID to be set after creation")
		}
	})

	t.Run("returns error for duplicate email", func(t *testing.T) {
		user1 := &TestUser{
			Name:  "User 1",
			Email: "duplicate@example.com",
			Age:   25,
		}
		repo.Create(ctx, user1)

		user2 := &TestUser{
			Name:  "User 2",
			Email: "duplicate@example.com", // Same email
			Age:   28,
		}

		err := repo.Create(ctx, user2)
		if err == nil {
			t.Error("Expected error for duplicate email")
		}
	})
}

func TestFindByID(t *testing.T) {
	db := setupTestDB(t)
	repo := New[TestUser](db)
	ctx := context.Background()

	// Create a test user
	createdUser := &TestUser{
		Name:  "Jane Doe",
		Email: "jane@example.com",
		Age:   28,
	}
	repo.Create(ctx, createdUser)

	t.Run("finds existing record by ID", func(t *testing.T) {
		var foundUser TestUser
		err := repo.FindByID(ctx, createdUser.ID, &foundUser)

		if err != nil {
			t.Fatalf("Failed to find user: %v", err)
		}
		if foundUser.ID != createdUser.ID {
			t.Errorf("Expected ID %d, got %d", createdUser.ID, foundUser.ID)
		}
		if foundUser.Name != "Jane Doe" {
			t.Errorf("Expected name Jane Doe, got %s", foundUser.Name)
		}
	})

	t.Run("returns error for non-existent ID", func(t *testing.T) {
		var user TestUser
		err := repo.FindByID(ctx, 99999, &user)

		if err == nil {
			t.Error("Expected error for non-existent ID")
		}
	})
}

func TestFindAll(t *testing.T) {
	db := setupTestDB(t)
	repo := New[TestUser](db)
	ctx := context.Background()

	// Create test users
	users := []TestUser{
		{Name: "User 1", Email: "user1@example.com", Age: 25},
		{Name: "User 2", Email: "user2@example.com", Age: 30},
		{Name: "User 3", Email: "user3@example.com", Age: 35},
	}
	for i := range users {
		repo.Create(ctx, &users[i])
	}

	t.Run("finds all records", func(t *testing.T) {
		foundUsers, err := repo.FindAll(ctx)

		if err != nil {
			t.Fatalf("Failed to find all users: %v", err)
		}
		if len(foundUsers) != 3 {
			t.Errorf("Expected 3 users, got %d", len(foundUsers))
		}
	})
}

func TestUpdate(t *testing.T) {
	db := setupTestDB(t)
	repo := New[TestUser](db)
	ctx := context.Background()

	// Create a test user
	user := &TestUser{
		Name:  "Original Name",
		Email: "original@example.com",
		Age:   25,
	}
	repo.Create(ctx, user)

	t.Run("updates record successfully", func(t *testing.T) {
		user.Name = "Updated Name"
		user.Age = 30

		err := repo.Update(ctx, user)
		if err != nil {
			t.Fatalf("Failed to update user: %v", err)
		}

		// Verify the update
		var updatedUser TestUser
		repo.FindByID(ctx, user.ID, &updatedUser)
		if updatedUser.Name != "Updated Name" {
			t.Errorf("Expected name to be updated to 'Updated Name', got %s", updatedUser.Name)
		}
		if updatedUser.Age != 30 {
			t.Errorf("Expected age to be updated to 30, got %d", updatedUser.Age)
		}
	})
}

func TestDelete(t *testing.T) {
	db := setupTestDB(t)
	repo := New[TestUser](db)
	ctx := context.Background()

	t.Run("deletes record successfully", func(t *testing.T) {
		user := &TestUser{
			Name:  "To Delete",
			Email: "delete@example.com",
			Age:   25,
		}
		repo.Create(ctx, user)

		err := repo.Delete(ctx, user)
		if err != nil {
			t.Fatalf("Failed to delete user: %v", err)
		}

		// Verify deletion
		var deletedUser TestUser
		err = repo.FindByID(ctx, user.ID, &deletedUser)
		if err == nil {
			t.Error("Expected error when finding deleted user")
		}
	})
}

func TestDeleteByID(t *testing.T) {
	db := setupTestDB(t)
	repo := New[TestUser](db)
	ctx := context.Background()

	t.Run("deletes record by ID successfully", func(t *testing.T) {
		user := &TestUser{
			Name:  "To Delete By ID",
			Email: "deletebyid@example.com",
			Age:   25,
		}
		repo.Create(ctx, user)

		err := repo.DeleteByID(ctx, user.ID)
		if err != nil {
			t.Fatalf("Failed to delete user by ID: %v", err)
		}

		// Verify deletion
		var deletedUser TestUser
		err = repo.FindByID(ctx, user.ID, &deletedUser)
		if err == nil {
			t.Error("Expected error when finding deleted user")
		}
	})
}

func TestCount(t *testing.T) {
	db := setupTestDB(t)
	repo := New[TestUser](db)
	ctx := context.Background()

	// Create test users
	for i := 1; i <= 5; i++ {
		user := &TestUser{
			Name:  "User",
			Email: "user" + string(rune(i)) + "@example.com",
			Age:   20 + i,
		}
		repo.Create(ctx, user)
	}

	t.Run("counts all records", func(t *testing.T) {
		count, err := repo.Count(ctx)
		if err != nil {
			t.Fatalf("Failed to count users: %v", err)
		}
		if count != 5 {
			t.Errorf("Expected count 5, got %d", count)
		}
	})
}

func TestFindWhere(t *testing.T) {
	db := setupTestDB(t)
	repo := New[TestUser](db)
	ctx := context.Background()

	// Create test users
	users := []TestUser{
		{Name: "Alice", Email: "alice@example.com", Age: 25},
		{Name: "Bob", Email: "bob@example.com", Age: 30},
		{Name: "Charlie", Email: "charlie@example.com", Age: 25},
	}
	for i := range users {
		repo.Create(ctx, &users[i])
	}

	t.Run("finds records matching condition", func(t *testing.T) {
		foundUsers, err := repo.FindWhere(ctx, "age = ?", 25)

		if err != nil {
			t.Fatalf("Failed to find users: %v", err)
		}
		if len(foundUsers) != 2 {
			t.Errorf("Expected 2 users with age 25, got %d", len(foundUsers))
		}
	})

	t.Run("returns empty slice when no match", func(t *testing.T) {
		foundUsers, err := repo.FindWhere(ctx, "age = ?", 99)

		if err != nil {
			t.Fatalf("Failed to find users: %v", err)
		}
		if len(foundUsers) != 0 {
			t.Errorf("Expected 0 users, got %d", len(foundUsers))
		}
	})
}

func TestFirstWhere(t *testing.T) {
	db := setupTestDB(t)
	repo := New[TestUser](db)
	ctx := context.Background()

	// Create test users
	user := &TestUser{
		Name:  "Test User",
		Email: "testuser@example.com",
		Age:   30,
	}
	repo.Create(ctx, user)

	t.Run("finds first record matching condition", func(t *testing.T) {
		var foundUser TestUser
		err := repo.FirstWhere(ctx, &foundUser, "email = ?", "testuser@example.com")

		if err != nil {
			t.Fatalf("Failed to find user: %v", err)
		}
		if foundUser.Email != "testuser@example.com" {
			t.Errorf("Expected email testuser@example.com, got %s", foundUser.Email)
		}
	})

	t.Run("returns error when no match", func(t *testing.T) {
		var foundUser TestUser
		err := repo.FirstWhere(ctx, &foundUser, "email = ?", "nonexistent@example.com")

		if err == nil {
			t.Error("Expected error when no user found")
		}
	})
}

func TestPaginate(t *testing.T) {
	db := setupTestDB(t)
	repo := New[TestUser](db)
	ctx := context.Background()

	// Create 15 test users
	for i := 1; i <= 15; i++ {
		user := &TestUser{
			Name:  "User",
			Email: "user" + string(rune(i)) + "@example.com",
			Age:   20 + i,
		}
		repo.Create(ctx, user)
	}

	t.Run("returns first page", func(t *testing.T) {
		users, total, err := repo.Paginate(ctx, 1, 5)

		if err != nil {
			t.Fatalf("Failed to paginate: %v", err)
		}
		if len(users) != 5 {
			t.Errorf("Expected 5 users on first page, got %d", len(users))
		}
		if total != 15 {
			t.Errorf("Expected total 15, got %d", total)
		}
	})

	t.Run("returns second page", func(t *testing.T) {
		users, total, err := repo.Paginate(ctx, 2, 5)

		if err != nil {
			t.Fatalf("Failed to paginate: %v", err)
		}
		if len(users) != 5 {
			t.Errorf("Expected 5 users on second page, got %d", len(users))
		}
		if total != 15 {
			t.Errorf("Expected total 15, got %d", total)
		}
	})

	t.Run("returns last page with remaining items", func(t *testing.T) {
		users, total, err := repo.Paginate(ctx, 3, 5)

		if err != nil {
			t.Fatalf("Failed to paginate: %v", err)
		}
		if len(users) != 5 {
			t.Errorf("Expected 5 users on third page, got %d", len(users))
		}
		if total != 15 {
			t.Errorf("Expected total 15, got %d", total)
		}
	})
}

func TestTransaction(t *testing.T) {
	db := setupTestDB(t)
	repo := New[TestUser](db)
	ctx := context.Background()

	t.Run("commits transaction on success", func(t *testing.T) {
		err := repo.Transaction(ctx, func(tx *gorm.DB) error {
			user1 := &TestUser{Name: "User 1", Email: "tx1@example.com", Age: 25}
			user2 := &TestUser{Name: "User 2", Email: "tx2@example.com", Age: 30}

			if err := tx.Create(user1).Error; err != nil {
				return err
			}
			if err := tx.Create(user2).Error; err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			t.Fatalf("Transaction failed: %v", err)
		}

		// Verify both users were created
		users, _ := repo.FindWhere(ctx, "email LIKE ?", "tx%@example.com")
		if len(users) != 2 {
			t.Errorf("Expected 2 users to be committed, got %d", len(users))
		}
	})

	t.Run("rolls back transaction on error", func(t *testing.T) {
		countBefore, _ := repo.Count(ctx)

		err := repo.Transaction(ctx, func(tx *gorm.DB) error {
			user := &TestUser{Name: "Rollback User", Email: "rollback@example.com", Age: 25}
			if err := tx.Create(user).Error; err != nil {
				return err
			}

			// Return an error to trigger rollback
			return gorm.ErrInvalidTransaction
		})

		if err == nil {
			t.Error("Expected transaction to fail")
		}

		// Verify no user was created
		countAfter, _ := repo.Count(ctx)
		if countAfter != countBefore {
			t.Errorf("Expected count to remain %d after rollback, got %d", countBefore, countAfter)
		}
	})
}
