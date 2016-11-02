package helpers

import "errors"

type MockedCaseDocumentRepository struct {
	Documents map[string]CaseDocument
}

func (r *MockedCaseDocumentRepository) FindCaseDocumentById(caseId string) (CaseDocument, error) {
	if doc, ok := r.Documents[caseId]; ok {
		return doc, nil
	} else {
		return doc, errors.New("Not found")
	}
}

func NewMockedCaseDocumentRepository() *MockedCaseDocumentRepository {
	return &MockedCaseDocumentRepository{make(map[string]CaseDocument)}
}

func NewMockedSingleCaseDocumentRepository(caseDocument CaseDocument) *MockedCaseDocumentRepository {
	repo := NewMockedCaseDocumentRepository()
	repo.Documents[caseDocument.Id.Hex()] = caseDocument
	return repo
}
