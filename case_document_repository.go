package helpers

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type CaseDocumentRepository interface {
	FindCaseDocumentById(caseId string) (CaseDocument, error)
}

func NewCaseDocumentRepository(mongoSession *mgo.Session, accountRepo AccountRepository) CaseDocumentRepository {
	return &mongoCaseDocumentRepository{mongoSession, accountRepo}
}

type mongoCaseDocumentRepository struct {
	mongoSession      *mgo.Session
	accountRepository AccountRepository
}

func (r *mongoCaseDocumentRepository) FindCaseDocumentById(caseId string) (CaseDocument, error) {
	var result CaseDocument

	err := r.mongoSession.DB(FindMongoDbString()).C("case_blocks.case_documents").FindId(bson.ObjectIdHex(caseId)).One(&result)
	if err != nil {
		return result, err
	}

	account, err := r.accountRepository.FindById(result.AccountId)
	if err != nil {
		return result, err
	}
	result.AccountCode = account.Nickname

	return result, nil
}
