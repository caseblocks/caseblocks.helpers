package helpers

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type CaseDocumentRepository interface {
	FindCaseDocumentById(caseId string) (CaseDocument, error)
}

func NewCaseDocumentRepository(mongoSession *mgo.Session) CaseDocumentRepository {
	return &mongoCaseDocumentRepository{mongoSession}
}

type mongoCaseDocumentRepository struct {
	mongoSession *mgo.Session
}

func (r *mongoCaseDocumentRepository) FindCaseDocumentById(caseId string) (CaseDocument, error) {
	var result CaseDocument
	err := r.mongoSession.DB(FindMongoDbString()).C("case_blocks.case_documents").FindId(bson.ObjectIdHex(caseId)).One(&result)
	return result, err
}
