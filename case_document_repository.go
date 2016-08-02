package helpers

import (
	"fmt"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type CaseDocumentRepository interface {
	FindCaseDocumentById(caseId string) (CaseDocument, error)
}

func NewCaseDocumentRepository(mongoSession *mgo.Session, accountRepo AccountRepository, log Log) CaseDocumentRepository {
	return &mongoCaseDocumentRepository{mongoSession, accountRepo, log}
}

type mongoCaseDocumentRepository struct {
	mongoSession      *mgo.Session
	accountRepository AccountRepository
	log               Log
}

func (r *mongoCaseDocumentRepository) FindCaseDocumentById(caseId string) (CaseDocument, error) {
	var result CaseDocument

	dbConn := FindMongoDbString()
	r.log.Debug(fmt.Sprintf("Retrieving case document %s from %s", caseId, dbConn))
	err := r.mongoSession.DB(dbConn).C("case_blocks.case_documents").FindId(bson.ObjectIdHex(caseId)).One(&result)
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
